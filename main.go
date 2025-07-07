// Command-manager is the entry point for the FreeBSD Command Manager CLI.
package main

import (
	"FreeBSD-Command-manager/cmd"
)

// VERSION is the current version of the FreeBSD Command Manager.
var VERSION = "0.02" //nolint: gochecknoglobals

var (
	commit = "0" //nolint: gochecknoglobals,unused
	built  = "0" //nolint: gochecknoglobals,unused
)

func main() {
	cmd.Execute()
}
