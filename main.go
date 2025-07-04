package main

import (
	"FreeBSD-Command-manager/cmd"
)

var VERSION = "0.01" //nolint: gochecknoglobals

var commit string = "0" //nolint: gochecknoglobals
var built string = "0"  //nolint: gochecknoglobals

func main() {
	cmd.Execute()
}
