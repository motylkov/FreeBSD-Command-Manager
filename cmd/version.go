package cmd

import (
	"FreeBSD-Command-manager/internal"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version info variables.
var (
	Version = "0"       //nolint: gochecknoglobals
	Commit  = "0"       //nolint: gochecknoglobals
	Built   = "0"       //nolint: gochecknoglobals
	Date    = "unknown" //nolint: gochecknoglobals
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		err := internal.Output(map[string]interface{}{
			"version": Version,
			"commit":  Commit,
			"built":   Built,
			"date":    Date,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() { //nolint
	cmd.AddCommand(versionCmd)
}
