package cmd

import (
	"fmt"
	"os"

	"FreeBSD-Command-manager/internal"
	"FreeBSD-Command-manager/internal/network/bareos"

	"github.com/spf13/cobra"
)

var (
	ifName       string
	delIfaceName string
)

var (
	bridgeName    string
	delBridgeName string
)

var (
	vlanName, vlanParent string
	vlanID               int
	delVlanName          string
)

var (
	greName, greRemote, greLocal string
	delGreName                   string
)

var (
	vxlanName, vxlanRemote, vxlanLocal, vxlanGroup, vxlanDev string
	vxlanID                                                  int
	delVxlanName                                             string
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
}

// iface
var ifaceCmd = &cobra.Command{
	Use:   "iface",
	Short: "Create a generic interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.CreateInterface(ifName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"interface": ifName,
			"status":    "created",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var delIfaceCmd = &cobra.Command{
	Use:   "delete-iface",
	Short: "Delete a network interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.DeleteInterface(delIfaceName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"interface": delIfaceName,
			"status":    "deleted",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

// bridge
var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Create a bridge interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.CreateBridge(bridgeName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"bridge": bridgeName,
			"status": "created",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var delBridgeCmd = &cobra.Command{
	Use:   "delete-bridge",
	Short: "Delete a bridge interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.DeleteBridge(delBridgeName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"bridge": delBridgeName,
			"status": "deleted",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

// vlan
var vlanCmd = &cobra.Command{
	Use:   "vlan",
	Short: "Create VLAN interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.CreateVLAN(vlanName, vlanParent, vlanID); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"vlan":    vlanName,
			"parent":  vlanParent,
			"vlan_id": vlanID,
			"status":  "created",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var delVlanCmd = &cobra.Command{
	Use:   "delete-vlan",
	Short: "Delete a VLAN interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.DeleteVLAN(delVlanName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"vlan":   delVlanName,
			"status": "deleted",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

// gre
var greCmd = &cobra.Command{
	Use:   "gre",
	Short: "Create a GRE tunnel interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.CreateGRE(greName, greRemote, greLocal); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"gre":    greName,
			"local":  greLocal,
			"remote": greRemote,
			"status": "created",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var delGreCmd = &cobra.Command{
	Use:   "delete-gre",
	Short: "Delete a GRE tunnel interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.DeleteGRE(delGreName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"gre":    delGreName,
			"status": "deleted",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

// vxlan
var vxlanCmd = &cobra.Command{
	Use:   "vxlan",
	Short: "Create a VXLAN tunnel interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.CreateVXLAN(vxlanName, vxlanLocal, vxlanRemote, vxlanGroup, vxlanDev, vxlanID); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"vxlan":  vxlanName,
			"local":  vxlanLocal,
			"remote": vxlanRemote,
			"group":  vxlanGroup,
			"dev":    vxlanDev,
			"vni":    vxlanID,
			"status": "created",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var delVxlanCmd = &cobra.Command{
	Use:   "delete-vxlan",
	Short: "Delete a VXLAN tunnel interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		if err := manager.DeleteVXLAN(delVxlanName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"vxlan":  delVxlanName,
			"status": "deleted",
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all network interfaces",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		interfaces, err := manager.List()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"interfaces": interfaces,
			"count":      len(interfaces),
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

var networkInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about a network interface",
	Run: func(cmd *cobra.Command, args []string) {
		manager := bareos.DefaultManager()
		info, err := manager.GetInfo(ifName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"interface_info": info,
		}
		if err := internal.Output(result); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	// iface
	networkCmd.AddCommand(ifaceCmd)
	ifaceCmd.Flags().StringVar(&ifName, "name", "", "Interface name (required)")
	if err := ifaceCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	networkCmd.AddCommand(delIfaceCmd)
	delIfaceCmd.Flags().StringVar(&delIfaceName, "name", "", "Interface name (required)")
	if err := delIfaceCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	// bridge
	networkCmd.AddCommand(bridgeCmd)
	bridgeCmd.Flags().StringVar(&bridgeName, "name", "", "Bridge interface name (required)")
	if err := bridgeCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	networkCmd.AddCommand(delBridgeCmd)
	delBridgeCmd.Flags().StringVar(&delBridgeName, "name", "", "Bridge interface name (required)")
	if err := delBridgeCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	// VLAN
	networkCmd.AddCommand(vlanCmd)
	vlanCmd.Flags().StringVar(&vlanName, "name", "", "VLAN interface name (required)")
	vlanCmd.Flags().StringVar(&vlanParent, "parent", "", "Parent interface (required)")
	vlanCmd.Flags().IntVar(&vlanID, "id", 0, "VLAN ID (required)")
	if err := vlanCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	if err := vlanCmd.MarkFlagRequired("parent"); err != nil {
		panic(err)
	}
	if err := vlanCmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}

	networkCmd.AddCommand(delVlanCmd)
	delVlanCmd.Flags().StringVar(&delVlanName, "name", "", "VLAN interface name (required)")
	if err := delVlanCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	// GRE
	networkCmd.AddCommand(greCmd)
	greCmd.Flags().StringVar(&greName, "name", "", "GRE interface name (required)")
	greCmd.Flags().StringVar(&greLocal, "local", "", "Local address (required)")
	greCmd.Flags().StringVar(&greRemote, "remote", "", "Remote address (required)")
	if err := greCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	if err := greCmd.MarkFlagRequired("local"); err != nil {
		panic(err)
	}
	if err := greCmd.MarkFlagRequired("remote"); err != nil {
		panic(err)
	}

	networkCmd.AddCommand(delGreCmd)
	delGreCmd.Flags().StringVar(&delGreName, "name", "", "GRE interface name (required)")
	if err := delGreCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	// vxlan
	networkCmd.AddCommand(vxlanCmd)
	vxlanCmd.Flags().StringVar(&vxlanName, "name", "", "VXLAN interface name (required)")
	vxlanCmd.Flags().StringVar(&vxlanLocal, "local", "", "Local address (required)")
	vxlanCmd.Flags().StringVar(&vxlanRemote, "remote", "", "Remote address (required)")
	vxlanCmd.Flags().StringVar(&vxlanGroup, "group", "", "VXLAN group")
	vxlanCmd.Flags().StringVar(&vxlanDev, "dev", "", "VXLAN device")
	vxlanCmd.Flags().IntVar(&vxlanID, "vni", 0, "VXLAN Network Identifier (required)")
	if err := vxlanCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	if err := vxlanCmd.MarkFlagRequired("local"); err != nil {
		panic(err)
	}
	if err := vxlanCmd.MarkFlagRequired("remote"); err != nil {
		panic(err)
	}
	if err := vxlanCmd.MarkFlagRequired("vni"); err != nil {
		panic(err)
	}

	networkCmd.AddCommand(delVxlanCmd)
	delVxlanCmd.Flags().StringVar(&delVxlanName, "name", "", "VXLAN interface name (required)")
	if err := delVxlanCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	// List and Info commands
	networkCmd.AddCommand(networkListCmd)

	networkCmd.AddCommand(networkInfoCmd)
	networkInfoCmd.Flags().StringVar(&ifName, "name", "", "Interface name (required)")
	if err := networkInfoCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	cmd.AddCommand(networkCmd)
}
