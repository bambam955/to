// Package main implements the to-backend binary — the Go backend for the
// 'to' directory navigation tool. The shell wrapper (wrappers/to.bash) invokes
// this binary and interprets its stdout protocol responses to perform cd.
package main

func main() {
	Execute()
}
