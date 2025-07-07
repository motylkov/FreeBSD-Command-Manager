// Package cmd provides the CLI commands for the FreeBSD Command Manager.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "manager",
	Short: "FreeBSD Manager CLI",
	Long:  `A CLI tool for managing FreeBSD`,
}

// Execute runs the root command for the FreeBSD Command Manager CLI.
func Execute(v, c, b, d string) {
	Version = v
	Commit = c
	Built = b
	Date = d
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
