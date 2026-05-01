#!/bin/sh

set -eu

REPO_OWNER="bambam955"
REPO_NAME="to"
API_ROOT="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}"
DOWNLOAD_ROOT="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download"

log() {
    printf '%s\n' "$*"
}

fail() {
    printf 'error: %s\n' "$*" >&2
    exit 1
}

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        fail "required command not found: $1"
    fi
}

extract_tag_name() {
    # The installer only needs the release tag_name field, so keep parsing
    # deliberately narrow instead of introducing a JSON parser dependency.
    printf '%s' "$1" |
        tr -d '\n' |
        sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p'
}

fetch_release_json() {
    url="$1"

    if ! curl -fsSL \
        -H "Accept: application/vnd.github+json" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${url}"; then
        return 1
    fi
}

resolve_version() {
    pinned_version="${TO_INSTALL_VERSION:-}"

    if [ -n "${pinned_version}" ]; then
        case "${pinned_version}" in
        *[!0-9.]*)
            fail "TO_INSTALL_VERSION must be a bare semantic version like 1.2.3"
            ;;
        *)
            ;;
        esac

        set +e
        release_json="$(fetch_release_json "${API_ROOT}/releases/tags/${pinned_version}")"
        release_status=$?
        set -e

        if [ "${release_status}" -ne 0 ]; then
            fail "failed to resolve GitHub release for TO_INSTALL_VERSION=${pinned_version}"
        fi
    else
        set +e
        release_json="$(fetch_release_json "${API_ROOT}/releases/latest")"
        release_status=$?
        set -e

        if [ "${release_status}" -ne 0 ]; then
            fail "failed to resolve the latest GitHub release"
        fi
    fi

    version="$(extract_tag_name "${release_json}")"
    if [ -z "${version}" ]; then
        fail "could not determine a release version from the GitHub API response"
    fi

    if [ -n "${pinned_version}" ] && [ "${version}" != "${pinned_version}" ]; then
        fail "resolved release ${version} does not match TO_INSTALL_VERSION=${pinned_version}"
    fi

    printf '%s\n' "${version}"
}

detect_os() {
    os_name="$(uname -s)"

    case "${os_name}" in
    Linux)
        printf 'linux\n'
        ;;
    *)
        fail "unsupported operating system: ${os_name} (Linux only for now)"
        ;;
    esac
}

detect_arch() {
    arch_name="$(uname -m)"

    case "${arch_name}" in
    x86_64)
        printf 'amd64\n'
        ;;
    aarch64 | arm64)
        printf 'arm64\n'
        ;;
    *)
        fail "unsupported architecture: ${arch_name} (expected x86_64, aarch64, or arm64)"
        ;;
    esac
}

resolve_install_dir() {
    requested_dir="${TO_INSTALL_DIR:-}"

    if [ -z "${requested_dir}" ]; then
        if [ -z "${HOME:-}" ]; then
            fail "HOME must be set when TO_INSTALL_DIR is not provided"
        fi
        requested_dir="${HOME}/.local/bin"
    fi

    case "${requested_dir}" in
    /*)
        absolute_dir="${requested_dir}"
        ;;
    *)
        absolute_dir="${PWD}/${requested_dir}"
        ;;
    esac

    if [ -e "${absolute_dir}" ] && [ ! -d "${absolute_dir}" ]; then
        fail "install path exists but is not a directory: ${absolute_dir}"
    fi

    if ! mkdir -p "${absolute_dir}"; then
        fail "failed to create install directory: ${absolute_dir}"
    fi

    # Canonicalize after creation so later copies and instructions use one
    # stable absolute path even when TO_INSTALL_DIR was relative.
    if ! (
        cd "${absolute_dir}" && pwd -P
    ); then
        fail "failed to resolve install directory: ${absolute_dir}"
    fi
}

escape_toml_basic_string() {
    printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

write_install_config() {
    install_dir="$1"
    if [ -z "${HOME:-}" ]; then
        return 0
    fi

    config_dir="${HOME}/.config/to"
    config_path="${config_dir}/config.toml"

    if ! mkdir -p "${config_dir}"; then
        fail "failed to create config directory: ${config_dir}"
    fi

    # Keep this field name and TOML shape aligned with pkg/config.Config.
    escaped_install_dir="$(escape_toml_basic_string "${install_dir}")"

    if ! printf 'install_dir = "%s"\n' "${escaped_install_dir}" >"${config_path}"; then
        fail "failed to write install config: ${config_path}"
    fi
}

compute_sha256() {
    file_path="$1"

    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "${file_path}" | awk '{ print $1 }'
        return 0
    fi

    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "${file_path}" | awk '{ print $1 }'
        return 0
    fi

    fail "required checksum tool not found: sha256sum or shasum"
}

download_file() {
    url="$1"
    destination="$2"

    if ! curl -fsSL "${url}" -o "${destination}"; then
        fail "failed to download ${url}"
    fi
}

path_contains() {
    case ":${PATH:-}:" in
    *:"$1":*)
        return 0
        ;;
    *)
        return 1
        ;;
    esac
}

print_shell_instructions() {
    install_dir="$1"
    shell_name="${2:-}"

    case "${shell_name}" in
    bash)
        log "Add this line to ~/.bashrc:"
        log "  source \"${install_dir}/to.bash\""
        ;;
    zsh)
        log "Add this line to ~/.zshrc:"
        log "  source \"${install_dir}/to.zsh\""
        ;;
    fish)
        log "Add this line to ~/.config/fish/config.fish:"
        log "  source \"${install_dir}/to.fish\""
        ;;
    *)
        log "Source the wrapper for your shell:"
        log "  bash: source \"${install_dir}/to.bash\""
        log "  zsh:  source \"${install_dir}/to.zsh\""
        log "  fish: source \"${install_dir}/to.fish\""
        ;;
    esac
}

require_command curl
require_command awk
require_command sed
require_command tar
require_command mktemp
require_command uname
require_command cp
require_command chmod
require_command mkdir

version="$(resolve_version)"
os="$(detect_os)"
arch="$(detect_arch)"
install_dir="$(resolve_install_dir)"
archive_name="to-${version}-${os}-${arch}.tar.gz"
checksums_name="checksums.txt"
archive_url="${DOWNLOAD_ROOT}/${version}/${archive_name}"
checksums_url="${DOWNLOAD_ROOT}/${version}/${checksums_name}"
shell_name="${SHELL:-}"
shell_name="${shell_name##*/}"

temp_dir="$(mktemp -d "${TMPDIR:-/tmp}/to-install.XXXXXX")"
cleanup() {
    rm -rf "${temp_dir}"
}
trap cleanup EXIT INT TERM HUP

archive_path="${temp_dir}/${archive_name}"
checksums_path="${temp_dir}/${checksums_name}"
extract_dir="${temp_dir}/extract"

log "Installing to ${install_dir}"
log "Resolved release ${version} for ${os}/${arch}"

download_file "${archive_url}" "${archive_path}"
download_file "${checksums_url}" "${checksums_path}"

expected_hash="$(awk -v file="${archive_name}" '$2 == file { print $1; exit }' "${checksums_path}")"
if [ -z "${expected_hash}" ]; then
    fail "checksum entry not found for ${archive_name}"
fi

actual_hash="$(compute_sha256 "${archive_path}")"
if [ "${actual_hash}" != "${expected_hash}" ]; then
    fail "checksum verification failed for ${archive_name}"
fi

mkdir -p "${extract_dir}"
tar -xzf "${archive_path}" -C "${extract_dir}"

for asset_name in to-backend to.bash to.zsh to.fish; do
    if [ ! -f "${extract_dir}/${asset_name}" ]; then
        fail "release archive is missing ${asset_name}"
    fi

    cp "${extract_dir}/${asset_name}" "${install_dir}/${asset_name}"
done

# Keep the installed files directly runnable and sourceable regardless of
# which shell archive permissions the caller's filesystem preserved.
chmod 0755 \
    "${install_dir}/to-backend" \
    "${install_dir}/to.bash" \
    "${install_dir}/to.zsh" \
    "${install_dir}/to.fish"

write_install_config "${install_dir}"

log ""
log "Installed to ${install_dir}"

set +e
path_contains "${install_dir}"
path_contains_status=$?
set -e

if [ "${path_contains_status}" -eq 0 ]; then
    log "PATH already contains ${install_dir}"
else
    log "Add ${install_dir} to your PATH if needed:"
    case "${shell_name}" in
    fish)
        log "  set -Ux fish_user_paths \"${install_dir}\" \$fish_user_paths"
        ;;
    *)
        log "  export PATH=\"${install_dir}:\$PATH\""
        ;;
    esac
fi

log ""
print_shell_instructions "${install_dir}" "${shell_name}"
