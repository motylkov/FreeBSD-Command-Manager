package ifconfig

import (
	"testing"
)

func TestParseIfconfig(t *testing.T) {
	// Test with real FreeBSD ifconfig output
	input := `em0: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	options=48505bb<RXCSUM,TXCSUM,VLAN_MTU,VLAN_HWTAGGING,JUMBO_MTU,VLAN_HWCSUM,TSO4,LRO,VLAN_HWFILTER,VLAN_HWTSO,HWSTATS,MEXTPG>
	ether 08:00:27:20:af:31
	inet 192.168.88.241 netmask 0xffffff00 broadcast 192.168.88.255
	inet6 fe80::a00:27ff:fe20:af31%em0 prefixlen 64 scopeid 0x1
	media: Ethernet autoselect (1000baseT <full-duplex>)
	status: active
	nd6 options=23<PERFORMNUD,ACCEPT_RTADV,AUTO_LINKLOCAL>
lo0: flags=1008049<UP,LOOPBACK,RUNNING,MULTICAST,LOWER_UP> metric 0 mtu 16384
	options=680003<RXCSUM,TXCSUM,LINKSTATE,RXCSUM_IPV6,TXCSUM_IPV6>
	inet 127.0.0.1 netmask 0xff000000
	inet6 ::1 prefixlen 128
	inet6 fe80::1%lo0 prefixlen 64 scopeid 0x2
	groups: lo
	nd6 options=21<PERFORMNUD,AUTO_LINKLOCAL>
salsa0: flags=1008842<BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	options=4000503<RXCSUM,TXCSUM,TSO4,LRO,MEXTPG>
	ether 08:00:27:20:af:31
	groups: vlan
	vlan: 100 vlanproto: 802.1q vlanpcp: 0 parent interface: em0
	media: Ethernet autoselect (1000baseT <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
vlan1: flags=8002<BROADCAST,MULTICAST> metric 0 mtu 1500
	options=0
	ether 00:00:00:00:00:00
	groups: vlan
	vlan: 0 vlanproto: 0x0000 vlanpcp: 0 parent interface: <none>
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>`

	interfaces := ParseIfconfig(input)

	// Should have 4 interfaces: em0, lo0, salsa0, vlan1
	if len(interfaces) != 4 {
		t.Errorf("expected 4 interfaces, got %d", len(interfaces))
	}

	// Test em0 interface
	em0 := findInterface(interfaces, "em0")
	if em0 == nil {
		t.Error("em0 interface not found")
	} else {
		if em0.Type != Ethernet {
			t.Errorf("expected em0 type to be %s, got %s", Ethernet, em0.Type)
		}
		if em0.Status != "up" {
			t.Errorf("expected em0 status to be up, got %s", em0.Status)
		}
		if len(em0.IPv4) != 1 || em0.IPv4[0] != "192.168.88.241/24" {
			t.Errorf("expected em0 IPv4 to be [192.168.88.241/24], got %v", em0.IPv4)
		}
		if len(em0.IPv6) != 1 || em0.IPv6[0] != "fe80::a00:27ff:fe20:af31/64" {
			t.Errorf("expected em0 IPv6 to be [fe80::a00:27ff:fe20:af31/64], got %v", em0.IPv6)
		}
		if em0.MAC != "08:00:27:20:af:31" {
			t.Errorf("expected em0 MAC to be 08:00:27:20:af:31, got %s", em0.MAC)
		}
	}

	// Test lo0 interface
	lo0 := findInterface(interfaces, "lo0")
	if lo0 == nil {
		t.Error("lo0 interface not found")
	} else {
		if lo0.Type != Loopback {
			t.Errorf("expected lo0 type to be %s, got %s", Loopback, lo0.Type)
		}
		if lo0.Status != "up" {
			t.Errorf("expected lo0 status to be up, got %s", lo0.Status)
		}
		if len(lo0.IPv4) != 1 || lo0.IPv4[0] != "127.0.0.1/8" {
			t.Errorf("expected lo0 IPv4 to be [127.0.0.1/8], got %v", lo0.IPv4)
		}
		if len(lo0.IPv6) != 2 {
			t.Errorf("expected lo0 to have 2 IPv6 addresses, got %d", len(lo0.IPv6))
		}
		// Check for specific IPv6 addresses with CIDR
		found128 := false
		found64 := false
		for _, ipv6 := range lo0.IPv6 {
			if ipv6 == "::1/128" {
				found128 = true
			}
			if ipv6 == "fe80::1/64" {
				found64 = true
			}
		}
		if !found128 {
			t.Error("expected lo0 to have ::1/128 IPv6 address")
		}
		if !found64 {
			t.Error("expected lo0 to have fe80::1/64 IPv6 address")
		}
	}

	// Test salsa0 interface (VLAN)
	salsa0 := findInterface(interfaces, "salsa0")
	if salsa0 == nil {
		t.Error("salsa0 interface not found")
	} else {
		if salsa0.Type != VLAN {
			t.Errorf("expected salsa0 type to be %s, got %s", VLAN, salsa0.Type)
		}
		if salsa0.Status != "up" {
			t.Errorf("expected salsa0 status to be up, got %s", salsa0.Status)
		}
		if len(salsa0.IPv4) != 0 {
			t.Errorf("expected salsa0 to have no IPv4 addresses, got %v", salsa0.IPv4)
		}
		if salsa0.MAC != "08:00:27:20:af:31" {
			t.Errorf("expected salsa0 MAC to be 08:00:27:20:af:31, got %s", salsa0.MAC)
		}
	}

	// Test vlan1 interface
	vlan1 := findInterface(interfaces, "vlan1")
	if vlan1 == nil {
		t.Error("vlan1 interface not found")
	} else {
		if vlan1.Type != VLAN {
			t.Errorf("expected vlan1 type to be %s, got %s", VLAN, vlan1.Type)
		}
		if vlan1.Status != "down" {
			t.Errorf("expected vlan1 status to be down, got %s", vlan1.Status)
		}
		if vlan1.MAC != "00:00:00:00:00:00" {
			t.Errorf("expected vlan1 MAC to be 00:00:00:00:00:00, got %s", vlan1.MAC)
		}
	}
}

func findInterface(interfaces []Info, name string) *Info {
	for _, iface := range interfaces {
		if iface.Name == name {
			return &iface
		}
	}
	return nil
}
