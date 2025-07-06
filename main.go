package main

import (
	"FreeBSD-Command-manager/cmd"
)

var VERSION = "0.01" //nolint: gochecknoglobals

var (
	commit = "0" //nolint: gochecknoglobals,unused
	built  = "0" //nolint: gochecknoglobals,unused
)

func main() {
	cmd.Execute()
}
