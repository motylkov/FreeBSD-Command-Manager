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
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := jail.DefaultManager()

		cfg := jail.Config{
			Name:  jailName,
			Path:  jailPath,
			IP:    jailIP,
			Mount: jailMount,
		}

		err := manager.Create(cfg)
		if err != nil {
			if e := internal.Output(map[string]interface{}{
				"error": err.Error(),
			}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if e := internal.Output(map[string]interface{}{
			"jail_id": jailName,
			"status":  "created",
			"ip":      jailIP,
			"network": "br-jails", // todo:  make dynamic !
		}); e != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var jailStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a jail",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := jail.DefaultManager()

		if err := manager.Start(jailName); err != nil {
			if e := internal.Output(map[string]interface{}{
				"error": err.Error(),
			}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}

		if err := internal.Output(map[string]interface{}{
			"jail_id": jailName,
			"status":  "started",
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var jailStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a jail",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := jail.DefaultManager()

		if err := manager.Stop(jailName); err != nil {
			if e := internal.Output(map[string]interface{}{
				"error": err.Error(),
			}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}

		if err := internal.Output(map[string]interface{}{
			"jail_id": jailName,
			"status":  "stopped",
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var jailDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a jail",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := jail.DefaultManager()

		if err := manager.Destroy(jailName); err != nil {
			if e := internal.Output(map[string]interface{}{
				"error": err.Error(),
			}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}

		if err := internal.Output(map[string]interface{}{
			"jail_id": jailName,
			"status":  "destroyed",
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var jailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all jails",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := jail.DefaultManager()

		jails, err := manager.List()
		if err != nil {
			if e := internal.Output(map[string]interface{}{
				"error": err.Error(),
			}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}

		if err := internal.Output(map[string]interface{}{
			"jails": jails,
			"count": len(jails),
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var jailInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about a jail",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := jail.DefaultManager()

		info, err := manager.GetInfo(jailName)
		if err != nil {
			if e := internal.Output(map[string]interface{}{
				"error": err.Error(),
			}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}

		if err := internal.Output(map[string]interface{}{
			"jail_info": info,
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() { //nolint
	// Create command flags
	jailCreateCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	jailCreateCmd.Flags().StringVar(&jailPath, "path", "", "Jail path (required)")
	jailCreateCmd.Flags().StringVar(&jailIP, "ip", "", "Jail IP address (required)")
	jailCreateCmd.Flags().StringVar(&jailMount, "mount", "", "ZFS dataset or image to mount (optional)")
	// check required params
	if err := jailCreateCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := jailCreateCmd.MarkFlagRequired("path"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := jailCreateCmd.MarkFlagRequired("ip"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Start command flags
	jailStartCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	// check required params
	if err := jailStartCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Stop command flags
	jailStopCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	// check required params
	if err := jailStopCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Destroy command flags
	jailDestroyCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	// check required params
	if err := jailDestroyCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Info command flags
	jailInfoCmd.Flags().StringVar(&jailName, "name", "", "Jail name (required)")
	// check required params
	if err := jailInfoCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Add all subcommands
	jailCmd.AddCommand(jailCreateCmd)
	jailCmd.AddCommand(jailStartCmd)
	jailCmd.AddCommand(jailStopCmd)
	jailCmd.AddCommand(jailDestroyCmd)
	jailCmd.AddCommand(jailListCmd)
	jailCmd.AddCommand(jailInfoCmd)

	cmd.AddCommand(jailCmd)
}
