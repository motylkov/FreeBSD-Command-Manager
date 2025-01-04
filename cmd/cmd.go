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

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
