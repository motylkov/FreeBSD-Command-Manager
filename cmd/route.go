package cmd

import (
	"FreeBSD-Command-manager/internal/network/bareos"
	"fmt"
	"os"

	"FreeBSD-Command-manager/internal"

	"github.com/spf13/cobra"
)

var (
	routeFamily string
	routeNet    string
	routeGW     string
	routeIface  string // new flag for interface
)

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage routes (add, del, list)",
}

var routeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a route",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		if routeNet == "" || routeGW == "" {
			if e := internal.Output(map[string]interface{}{"error": "--net and --gw are required"}); e != nil {
				fmt.Fprintln(os.Stderr, e)
				os.Exit(1)
			}
			return
		}
		err := bareos.AddRouteWithIface(routeFamily, routeNet, routeGW, routeIface)
		if err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, e)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"status": "Route added successfully"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var routeDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a route",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		if routeNet == "" {
			if e := internal.Output(map[string]interface{}{"error": "--net is required"}); e != nil {
				fmt.Fprintln(os.Stderr, e)
				os.Exit(1)
			}
			return
		}
		err := bareos.DelRoute(routeFamily, routeNet)
		if err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, e)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"status": "Route deleted successfully"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var routeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List routes",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		out, err := bareos.ListAllRoutes(routeFamily)
		if err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, e)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"routes": out, "count": len(out)}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() { //nolint
	routeAddCmd.Flags().StringVar(&routeFamily, "family", "inet", "Address family (inet or inet6)")
	routeAddCmd.Flags().StringVar(&routeNet, "net", "", "Network or destination (required)")
	routeAddCmd.Flags().StringVar(&routeGW, "gw", "", "Gateway (required)")
	routeAddCmd.Flags().StringVar(&routeIface, "iface", "", "Outgoing interface (optional)")

	routeDelCmd.Flags().StringVar(&routeFamily, "family", "inet", "Address family (inet or inet6)")
	routeDelCmd.Flags().StringVar(&routeNet, "net", "", "Network or destination (required)")

	routeListCmd.Flags().StringVar(&routeFamily, "family", "", "Address family (inet, inet6, or empty for all)")

	routeCmd.AddCommand(routeAddCmd)
	routeCmd.AddCommand(routeDelCmd)
	routeCmd.AddCommand(routeListCmd)

	// Add route to top level, do not delete.
	// cmd.AddCommand(routeCmd)
}
