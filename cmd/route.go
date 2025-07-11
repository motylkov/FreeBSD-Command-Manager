package cmd

import (
	"FreeBSD-Command-manager/internal/network/bareos"
	"fmt"
	"os"

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
	Run: func(_ *cobra.Command, _ []string) {
		if routeNet == "" || routeGW == "" {
			fmt.Fprintln(os.Stderr, "--net and --gw are required")
			os.Exit(1)
		}
		err := bareos.AddRouteWithIface(routeFamily, routeNet, routeGW, routeIface)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Route added successfully")
	},
}

var routeDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a route",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		if routeNet == "" {
			fmt.Fprintln(os.Stderr, "--net is required")
			os.Exit(1)
		}
		err := bareos.DelRoute(routeFamily, routeNet)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Route deleted successfully")
	},
}

var routeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List routes",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		out, err := bareos.ListAllRoutes(routeFamily)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Print(out)
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

	cmd.AddCommand(routeCmd)
}
