package cmd

import (
	"fmt"

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
		fmt.Printf("%s.%s.%s.%s\n", Version, Commit, Built, Date)
	},
}

func init() { //nolint
	cmd.AddCommand(versionCmd)
}
