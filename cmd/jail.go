package cmd

import (
	"FreeBSD-Command-manager/internal"
	"FreeBSD-Command-manager/internal/jail"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jailName, jailPath, jailIP, jailMount string

var jailCmd = &cobra.Command{
	Use:   "jail",
	Short: "Manage FreeBSD jails",
}

var jailCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new jail",
	Run: func(cmd *cobra.Command, args []string) {
		manager := jail.DefaultManager()

		cfg := jail.Config{
			Name:  jailName,
			Path:  jailPath,
			IP:    jailIP,
			Mount: jailMount,
		}

		if err := manager.Create(cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		result := map[string]interface{}{
			"jail_id": jailName,
			"status":  "created",
			"ip":      jailIP,
			"network": "br-jails", // make dynamic ?
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var jailStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a jail",
	Run: func(cmd *cobra.Command, args []string) {
		manager := jail.DefaultManager()

		if err := manager.Start(jailName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		result := map[string]interface{}{
			"jail_id": jailName,
			"status":  "started",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var jailStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a jail",
	Run: func(cmd *cobra.Command, args []string) {
		manager := jail.DefaultManager()

		if err := manager.Stop(jailName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		result := map[string]interface{}{
			"jail_id": jailName,
			"status":  "stopped",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var jailDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a jail",
	Run: func(cmd *cobra.Command, args []string) {
		manager := jail.DefaultManager()

		if err := manager.Destroy(jailName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		result := map[string]interface{}{
			"jail_id": jailName,
			"status":  "destroyed",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var jailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all jails",
	Run: func(cmd *cobra.Command, args []string) {
		manager := jail.DefaultManager()

		jails, err := manager.List()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		result := map[string]interface{}{
			"jails": jails,
			"count": len(jails),
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var jailInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about a jail",
	Run: func(cmd *cobra.Command, args []string) {
		manager := jail.DefaultManager()

		info, err := manager.GetInfo(jailName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		result := map[string]interface{}{
			"jail_info": info,
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	// Create command flags
	jailCreateCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	jailCreateCmd.Flags().StringVar(&jailPath, "path", "", "Jail path (required)")
	jailCreateCmd.Flags().StringVar(&jailIP, "ip", "", "Jail IP address (required)")
	jailCreateCmd.Flags().StringVar(&jailMount, "mount", "", "ZFS dataset or image to mount (optional)")
	jailCreateCmd.MarkFlagRequired("name")
	jailCreateCmd.MarkFlagRequired("path")
	jailCreateCmd.MarkFlagRequired("ip")

	// Start command flags
	jailStartCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	jailStartCmd.MarkFlagRequired("name")

	// Stop command flags
	jailStopCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	jailStopCmd.MarkFlagRequired("name")

	// Destroy command flags
	jailDestroyCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	jailDestroyCmd.MarkFlagRequired("name")

	// Info command flags
	jailInfoCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	jailInfoCmd.MarkFlagRequired("name")

	// Add all subcommands
	jailCmd.AddCommand(jailCreateCmd)
	jailCmd.AddCommand(jailStartCmd)
	jailCmd.AddCommand(jailStopCmd)
	jailCmd.AddCommand(jailDestroyCmd)
	jailCmd.AddCommand(jailListCmd)
	jailCmd.AddCommand(jailInfoCmd)

	cmd.AddCommand(jailCmd)
}
