package cmd

<<<<<<< HEAD
import "github.com/spf13/cobra"

var netCmd = &cobra.Command{
	Use:   "jail",
	Short: "Manage FreeBSD network",
}

// ip

// iface

// bridge

// vlan
=======
import (
	"FreeBSD-Command-manager/internal/network"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ifName string
var delIfaceName string

var bridgeName string
var delBridgeName string

var vlanName, vlanParent string
var vlanID int
var delVlanName string

var greName, greRemote, greLocal string
var delGreName string

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
}

// ip (in next time)

// iface
var ifaceCmd = &cobra.Command{
	Use:   "iface",
	Short: "Create a generic interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.CreateInterface(ifName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Interface %s created\n", ifName)
	},
}

var delIfaceCmd = &cobra.Command{
	Use:   "delete-iface",
	Short: "Delete a network interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.DeleteInterface(delIfaceName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Interface %s deleted\n", delIfaceName)
	},
}

// bridge
var bridgeCmd = &cobra.Command{
	Use:   "bridge-iface",
	Short: "Create a bridge interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.CreateBridge(bridgeName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Bridge interface %s created\n", bridgeName)
	},
}

var delBridgeCmd = &cobra.Command{
	Use:   "delete-bridge-iface",
	Short: "Delete a bridge interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.DeleteBridge(delBridgeName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Bridge interface %s deleted\n", delBridgeName)
	},
}

// vlan
var vlanCmd = &cobra.Command{
	Use:   "vlan",
	Short: "Create VLAN interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.CreateVLAN(vlanName, vlanParent, vlanID); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("VLAN %s created on %s with id %d\n", vlanName, vlanParent, vlanID)
	},
}

var delVlanCmd = &cobra.Command{
	Use:   "delete-vlan",
	Short: "Delete a VLAN interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.DeleteVLAN(delVlanName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Interface %s deleted\n", delVlanName)
	},
}
>>>>>>> test-jail

// vxlan

// gre
<<<<<<< HEAD

func init() {
	// flags and subcommands
=======
var greCmd = &cobra.Command{
	Use:   "gre",
	Short: "Create a GRE tunnel interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.CreateGRE(greName, greRemote, greLocal); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("GRE %s created: local %s remote %s\n", greName, greLocal, greRemote)
	},
}

var delGreCmd = &cobra.Command{
	Use:   "gre",
	Short: "Create a GRE tunnel interface",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := network.NewBareOSManager()
		if err := mgr.DeleteGRE(delGreName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("GRE %s created: local %s remote %s\n", delGreName, greLocal, greRemote)
	},
}

// --

func init() {
	// iface
	networkCmd.AddCommand(ifaceCmd)
	ifaceCmd.Flags().StringVar(&ifName, "name", "", "Interface name (required)")
	ifaceCmd.MarkFlagRequired("type")
	ifaceCmd.MarkFlagRequired("name")

	networkCmd.AddCommand(delIfaceCmd)
	delIfaceCmd.Flags().StringVar(&delIfaceName, "name", "", "Interface name (required)")
	delIfaceCmd.MarkFlagRequired("name")

	// bridge
	networkCmd.AddCommand(bridgeCmd)
	ifaceCmd.Flags().StringVar(&ifName, "name", "", "Bridge interface name (required)")
	ifaceCmd.MarkFlagRequired("name")

	networkCmd.AddCommand(delBridgeCmd)
	delBridgeCmd.Flags().StringVar(&ifName, "name", "", "Bridge interface name (required)")
	delBridgeCmd.MarkFlagRequired("name")
	// addToBridgeCmd
	// removeFromBridgeCmd

	// VLAN
	networkCmd.AddCommand(vlanCmd)
	vlanCmd.Flags().StringVar(&vlanName, "name", "", "VLAN interface name (required)")
	vlanCmd.Flags().StringVar(&vlanParent, "parent", "", "Parent interface (required)")
	vlanCmd.Flags().IntVar(&vlanID, "id", 0, "VLAN ID (required)")
	vlanCmd.MarkFlagRequired("name")
	vlanCmd.MarkFlagRequired("parent")
	vlanCmd.MarkFlagRequired("id")

	networkCmd.AddCommand(delVlanCmd)
	delVlanCmd.Flags().StringVar(&vlanName, "name", "", "VLAN interface name (required)")
	delVlanCmd.MarkFlagRequired("name")

	// GRE
	networkCmd.AddCommand(greCmd)
	greCmd.Flags().StringVar(&greName, "name", "", "GRE interface name (required)")
	greCmd.Flags().StringVar(&greLocal, "local", "", "Local address (required)")
	greCmd.Flags().StringVar(&greRemote, "remote", "", "Remote address (required)")
	greCmd.MarkFlagRequired("name")
	greCmd.MarkFlagRequired("local")
	greCmd.MarkFlagRequired("remote")

	networkCmd.AddCommand(delGreCmd)
	delGreCmd.Flags().StringVar(&greName, "name", "", "GRE interface name (required)")
	delGreCmd.MarkFlagRequired("name")
>>>>>>> test-jail
}
