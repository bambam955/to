package main

import (
	"strings"

	"github.com/spf13/cobra"
)

const developmentVersion = "dev"

// buildVersion is injected at build time with ldflags for tagged releases.
// Local and CI builds default to "dev" so the binary always reports a version.
var buildVersion = developmentVersion

const versionTemplate = "TO {{.Version}}\n"

// configureVersion stamps the root command with the build-time version and the
// machine-readable version banner used by Cobra's built-in version flag.
func configureVersion(cmd *cobra.Command) {
	cmd.Version = resolvedVersion()
	cmd.SetVersionTemplate(versionTemplate)
}

// resolvedVersion normalizes empty build metadata back to the development
// fallback so binaries built outside a tag still report a useful version.
func resolvedVersion() string {
	if strings.TrimSpace(buildVersion) == "" {
		return developmentVersion
	}
	return buildVersion
}
