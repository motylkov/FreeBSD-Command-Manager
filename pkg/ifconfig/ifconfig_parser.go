package ifconfig

import (
	"fmt"
	"math/bits"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	Ethernet   string = "ethernet"
	Wireless   string = "wireless"
	Loopback   string = "loopback"
	Bridge     string = "bridge"
	VLAN       string = "vlan"
	VXLAN      string = "vxlan"
	Tunnel     string = "tunnel"
	PPP        string = "ppp"
	LAGG       string = "lagg"
	GIF        string = "gif"
	GRE        string = "gre"
	Tap        string = "tap"
	Stf        string = "stf"
	Enc        string = "enc"
	Unknown    string = "unknown"
	StatusUp   string = "up"
	StatusDown string = "down"
)

// Info represents information about a network interface
type Info struct {
	Name   string
	Type   string
	Status string
	IPv4   []string
	IPv6   []string
	MAC    string
}

// Parses FreeBSD ifconfig output into a slice of Info structs.
func ParseIfconfig(output string) []Info {
	var result []Info
	lines := strings.Split(output, "\n")

	var currentInfo *Info
	var flags string
	var media string
	var groups string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if isNewInterfaceLine(line) {
			if currentInfo != nil {
				currentInfo.Type = determineInterfaceType(flags, media, groups)
				currentInfo.Status = determineInterfaceStatus(flags)
				result = append(result, *currentInfo)
			}
			currentInfo = &Info{Name: extractInterfaceName(line)}
			flags, media, groups = extractFlagsMediaGroups(line)
			continue
		}

		if currentInfo == nil {
			continue
		}

		flags, media, groups = updateFlagsMediaGroups(line, flags, media, groups)
		parseMAC(line, currentInfo)
		parseIPv4(line, currentInfo)
		parseIPv6(line, currentInfo)
	}

	if currentInfo != nil {
		currentInfo.Type = determineInterfaceType(flags, media, groups)
		currentInfo.Status = determineInterfaceStatus(flags)
		result = append(result, *currentInfo)
	}

	return result
}

func isNewInterfaceLine(line string) bool {
	return regexp.MustCompile(`^([a-z0-9]+):\s+flags=`).MatchString(line)
}

func extractInterfaceName(line string) string {
	matches := regexp.MustCompile(`^([a-z0-9]+):\s+flags=`).FindStringSubmatch(line)
	if matches != nil {
		return matches[1]
	}
	return ""
}

func extractFlagsMediaGroups(line string) (flags, media, groups string) {
	if strings.Contains(line, "flags=") {
		flags = line
	}
	return
}

func updateFlagsMediaGroups(line, flags, media, groups string) (string, string, string) {
	if flags == "" && strings.Contains(line, "flags=") {
		flags = line
	}
	if strings.Contains(line, "media:") {
		media = line
	}
	if strings.Contains(line, "groups:") {
		groups = line
	}
	return flags, media, groups
}

func parseMAC(line string, currentInfo *Info) {
	if strings.Contains(line, "ether ") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			if mac, err := net.ParseMAC(parts[1]); err == nil {
				currentInfo.MAC = mac.String()
			}
		}
	}
}

func parseIPv4(line string, currentInfo *Info) {
	if strings.Contains(line, "inet ") && !strings.Contains(line, "inet6") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			ip := net.ParseIP(parts[1])
			if ip != nil && ip.To4() != nil {
				cidr := parts[1]
				if len(parts) >= 3 && strings.Contains(parts[2], "netmask") {
					netmaskStr := parts[3]
					if netmaskStr != "" && strings.HasPrefix(netmaskStr, "0x") {
						if cidrBits := hexNetmaskToCIDR(netmaskStr); cidrBits > 0 {
							cidr = parts[1] + "/" + fmt.Sprintf("%d", cidrBits)
						}
					}
				}
				currentInfo.IPv4 = append(currentInfo.IPv4, cidr)
			}
		}
	}
}

func parseIPv6(line string, currentInfo *Info) {
	if strings.Contains(line, "inet6 ") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			ipStr := parts[1]
			if idx := strings.Index(ipStr, "%"); idx != -1 {
				ipStr = ipStr[:idx]
			}
			ip := net.ParseIP(ipStr)
			if ip != nil && ip.To16() != nil && ip.To4() == nil {
				cidr := ipStr
				if len(parts) >= 3 && strings.Contains(parts[2], "prefixlen") {
					prefixLen := parts[3]
					if prefixLen != "" {
						cidr = ipStr + "/" + prefixLen
					}
				}
				currentInfo.IPv6 = append(currentInfo.IPv6, cidr)
			}
		}
	}
}

func determineInterfaceType(flags, media, groups string) string {
	// Check for loopback
	if strings.Contains(flags, "LOOPBACK") {
		return Loopback
	}

	// Check for bridge
	if strings.Contains(media, "bridge members:") {
		return Bridge
	}

	// Check for VLAN - look for vlan: in media or groups: vlan in flags
	if strings.Contains(media, "vlan:") || strings.Contains(groups, "vlan") {
		return VLAN
	}

	// Check for VXLAN
	if strings.Contains(media, "vxlan:") {
		return VXLAN
	}

	// Check for LAGG
	if strings.Contains(media, "laggproto:") {
		return LAGG
	}

	// Check for PPP
	if strings.Contains(flags, "POINTOPOINT") && strings.Contains(flags, "MULTICAST") {
		return PPP
	}

	// Check for wireless
	if strings.Contains(media, "IEEE 802.11") {
		return Wireless
	}

	// Check for Ethernet
	if strings.Contains(media, "Ethernet") || strings.Contains(media, "autoselect") {
		return Ethernet
	}

	// Check for tunnel
	if strings.Contains(flags, "TUNNEL") {
		return Tunnel
	}

	// Check for GIF
	if strings.Contains(flags, "SIMPLEX") && strings.Contains(flags, "MULTICAST") {
		return GIF
	}

	// Check for GRE
	if strings.Contains(media, "gre") {
		return GRE
	}

	// Check for TAP
	if strings.Contains(flags, "TAP") {
		return Tap
	}

	// Check for STF
	if strings.Contains(flags, "stf") {
		return Stf
	}

	// Check for ENC
	if strings.Contains(flags, "ENCAP") {
		return Enc
	}

	return Unknown
}

func determineInterfaceStatus(flags string) string {
	if strings.Contains(flags, "RUNNING") {
		return StatusUp
	}
	return StatusDown
}

func hexNetmaskToCIDR(netmaskStr string) int {
	netmask, err := strconv.ParseUint(netmaskStr, 0, 64)
	if err != nil {
		return 0
	}

	cidrBits := 0
	for netmask > 0 {
		cidrBits += bits.OnesCount64(netmask & 0xFF)
		netmask >>= 8
	}

	return cidrBits
}
