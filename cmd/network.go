package cmd

import (
	"FreeBSD-Command-manager/internal"
	"FreeBSD-Command-manager/internal/network/bareos"
	"fmt"
	"os"

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
	bridgeInterfaceName     string
	bridgeInterfaceToAdd    string
	bridgeInterfaceToRemove string
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
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.CreateInterface(ifName); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"interface": ifName, "status": "created"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var delIfaceCmd = &cobra.Command{
	Use:   "delete-iface",
	Short: "Delete a network interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.DeleteInterface(delIfaceName); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"interface": delIfaceName, "status": "deleted"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// bridge
var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Create a bridge interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.CreateBridge(bridgeName); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"bridge": bridgeName, "status": "created"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var delBridgeCmd = &cobra.Command{
	Use:   "delete-bridge",
	Short: "Delete a bridge interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.DeleteBridge(delBridgeName); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"bridge": delBridgeName, "status": "deleted"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var addInterfaceToBridgeCmd = &cobra.Command{
	Use:   "add-interface-to-bridge",
	Short: "Add an interface to a bridge",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.AddInterfaceToBridge(bridgeInterfaceName, bridgeInterfaceToAdd); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"bridge": bridgeInterfaceName, "interface": bridgeInterfaceToAdd, "status": "added"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var removeInterfaceFromBridgeCmd = &cobra.Command{
	Use:   "remove-interface-from-bridge",
	Short: "Remove an interface from a bridge",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.RemoveInterfaceFromBridge(bridgeInterfaceName, bridgeInterfaceToRemove); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"bridge": bridgeInterfaceName, "interface": bridgeInterfaceToRemove, "status": "removed"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// vlan
var vlanCmd = &cobra.Command{
	Use:   "vlan",
	Short: "Create VLAN interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.CreateVLAN(vlanName, vlanParent, vlanID); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"vlan": vlanName, "parent": vlanParent, "vlan_id": vlanID, "status": "created"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var delVlanCmd = &cobra.Command{
	Use:   "delete-vlan",
	Short: "Delete a VLAN interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.DeleteVLAN(delVlanName); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"vlan": delVlanName, "status": "deleted"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// gre
var greCmd = &cobra.Command{
	Use:   "gre",
	Short: "Create a GRE tunnel interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.CreateGRE(greName, greRemote, greLocal); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"gre": greName, "remote": greRemote, "local": greLocal, "status": "created"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var delGreCmd = &cobra.Command{
	Use:   "delete-gre",
	Short: "Delete a GRE tunnel interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.DeleteGRE(delGreName); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"gre": delGreName, "status": "deleted"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// vxlan
var vxlanCmd = &cobra.Command{
	Use:   "vxlan",
	Short: "Create a VXLAN tunnel interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.CreateVXLAN(vxlanName, vxlanLocal, vxlanRemote, vxlanGroup, vxlanDev, vxlanID); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"vxlan": vxlanName, "local": vxlanLocal, "remote": vxlanRemote, "group": vxlanGroup, "dev": vxlanDev, "vni": vxlanID, "status": "created"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var delVxlanCmd = &cobra.Command{
	Use:   "delete-vxlan",
	Short: "Delete a VXLAN tunnel interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		if err := manager.DeleteVXLAN(delVxlanName); err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"vxlan": delVxlanName, "status": "deleted"}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all network interfaces",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		interfaces, err := manager.List()
		if err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"interfaces": interfaces, "count": len(interfaces)}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var networkInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about a network interface",
	Run: func(cmd *cobra.Command, args []string) { //nolint:revive // cmd is required by cobra interface
		manager := bareos.DefaultManager()
		info, err := manager.GetInfo(ifName)
		if err != nil {
			if e := internal.Output(map[string]interface{}{"error": err.Error()}); e != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		if err := internal.Output(map[string]interface{}{"interface_info": info}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() { //nolint
	// iface
	networkCmd.AddCommand(ifaceCmd)
	ifaceCmd.Flags().StringVar(&ifName, "name", "", "Interface name (required)")
	if err := ifaceCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(delIfaceCmd)
	delIfaceCmd.Flags().StringVar(&delIfaceName, "name", "", "Interface name (required)")
	if err := delIfaceCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// bridge
	networkCmd.AddCommand(bridgeCmd)
	bridgeCmd.Flags().StringVar(&bridgeName, "name", "", "Bridge interface name (required)")
	if err := bridgeCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(delBridgeCmd)
	delBridgeCmd.Flags().StringVar(&delBridgeName, "name", "", "Bridge interface name (required)")
	if err := delBridgeCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(addInterfaceToBridgeCmd)
	addInterfaceToBridgeCmd.Flags().StringVar(&bridgeInterfaceName, "bridge", "", "Bridge name (required)")
	addInterfaceToBridgeCmd.Flags().StringVar(&bridgeInterfaceToAdd, "interface", "", "Interface to add (required)")
	if err := addInterfaceToBridgeCmd.MarkFlagRequired("bridge"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := addInterfaceToBridgeCmd.MarkFlagRequired("interface"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(removeInterfaceFromBridgeCmd)
	removeInterfaceFromBridgeCmd.Flags().StringVar(&bridgeInterfaceName, "bridge", "", "Bridge name (required)")
	removeInterfaceFromBridgeCmd.Flags().StringVar(&bridgeInterfaceToRemove, "interface", "", "Interface to remove (required)")
	if err := removeInterfaceFromBridgeCmd.MarkFlagRequired("bridge"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := removeInterfaceFromBridgeCmd.MarkFlagRequired("interface"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// VLAN
	networkCmd.AddCommand(vlanCmd)
	vlanCmd.Flags().StringVar(&vlanName, "name", "", "VLAN interface name (required)")
	vlanCmd.Flags().StringVar(&vlanParent, "parent", "", "Parent interface (required)")
	vlanCmd.Flags().IntVar(&vlanID, "id", 0, "VLAN ID (required)")
	if err := vlanCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := vlanCmd.MarkFlagRequired("parent"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := vlanCmd.MarkFlagRequired("id"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(delVlanCmd)
	delVlanCmd.Flags().StringVar(&delVlanName, "name", "", "VLAN interface name (required)")
	if err := delVlanCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// GRE
	networkCmd.AddCommand(greCmd)
	greCmd.Flags().StringVar(&greName, "name", "", "GRE interface name (required)")
	greCmd.Flags().StringVar(&greLocal, "local", "", "Local address (required)")
	greCmd.Flags().StringVar(&greRemote, "remote", "", "Remote address (required)")
	if err := greCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := greCmd.MarkFlagRequired("local"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := greCmd.MarkFlagRequired("remote"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(delGreCmd)
	delGreCmd.Flags().StringVar(&delGreName, "name", "", "GRE interface name (required)")
	if err := delGreCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := vxlanCmd.MarkFlagRequired("local"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := vxlanCmd.MarkFlagRequired("remote"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := vxlanCmd.MarkFlagRequired("vni"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(delVxlanCmd)
	delVxlanCmd.Flags().StringVar(&delVxlanName, "name", "", "VXLAN interface name (required)")
	if err := delVxlanCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// List and Info commands
	networkCmd.AddCommand(networkListCmd)

	networkCmd.AddCommand(networkInfoCmd)
	networkInfoCmd.Flags().StringVar(&ifName, "name", "", "Interface name (required)")
	if err := networkInfoCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	networkCmd.AddCommand(routeCmd)

	cmd.AddCommand(networkCmd)
}
