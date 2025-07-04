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
	Ethernet string = "ethernet"
	Wireless string = "wireless"
	Loopback string = "loopback"
	Bridge   string = "bridge"
	VLAN     string = "vlan"
	VXLAN    string = "vxlan"
	Tunnel   string = "tunnel"
	PPP      string = "ppp"
	LAGG     string = "lagg"
	GIF      string = "gif"
	GRE      string = "gre"
	Tap      string = "tap"
	Stf      string = "stf"
	Enc      string = "enc"
	Unknown  string = "unknown"
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

		// Check for new interface - only match actual interface names
		// This regex matches interface names like em0, lo0, vlan1, etc.
		// but excludes lines like "media:", "status:", "groups:", etc.
		if matches := regexp.MustCompile(`^([a-z0-9]+):\s+flags=`).FindStringSubmatch(line); matches != nil {
			// Save previous interface info
			if currentInfo != nil {
				currentInfo.Type = determineInterfaceType(flags, media, groups)
				currentInfo.Status = determineInterfaceStatus(flags)
				result = append(result, *currentInfo)
			}

			// Start new interface
			currentInfo = &Info{Name: matches[1]}
			flags = ""
			media = ""
			groups = ""

			// Extract flags from the same line
			if strings.Contains(line, "flags=") {
				flags = line
			}
			continue
		}

		if currentInfo == nil {
			continue
		}

		// Collect flags (if not already collected from interface line)
		if flags == "" && strings.Contains(line, "flags=") {
			flags = line
		}

		// Collect media info
		if strings.Contains(line, "media:") {
			media = line
		}

		// Collect groups info
		if strings.Contains(line, "groups:") {
			groups = line
		}

		// Parse MAC address
		if strings.Contains(line, "ether ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if mac, err := net.ParseMAC(parts[1]); err == nil {
					currentInfo.MAC = mac.String()
				}
			}
		}

		// Parse IPv4 addresses
		if strings.Contains(line, "inet ") && !strings.Contains(line, "inet6") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				ip := net.ParseIP(parts[1])
				if ip != nil && ip.To4() != nil {
					// Extract netmask if available
					cidr := parts[1]
					if len(parts) >= 3 && strings.Contains(parts[2], "netmask") {
						netmaskStr := parts[3] // netmask is the 4th field after "netmask"
						if netmaskStr != "" {
							// Convert hex netmask to CIDR
							if strings.HasPrefix(netmaskStr, "0x") {
								if cidrBits := hexNetmaskToCIDR(netmaskStr); cidrBits > 0 {
									cidr = parts[1] + "/" + fmt.Sprintf("%d", cidrBits)
								}
							}
						}
					}
					currentInfo.IPv4 = append(currentInfo.IPv4, cidr)
				}
			}
		}

		// Parse IPv6 addresses
		if strings.Contains(line, "inet6 ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// Handle scope identifiers (e.g., fe80::1%lo0)
				ipStr := parts[1]
				if idx := strings.Index(ipStr, "%"); idx != -1 {
					ipStr = ipStr[:idx]
				}

				ip := net.ParseIP(ipStr)
				if ip != nil && ip.To16() != nil && ip.To4() == nil {
					// Extract prefix length if available
					cidr := ipStr
					if len(parts) >= 3 && strings.Contains(parts[2], "prefixlen") {
						prefixLen := parts[3] // prefixlen is the 4th field after "prefixlen"
						if prefixLen != "" {
							cidr = ipStr + "/" + prefixLen
						}
					}
					currentInfo.IPv6 = append(currentInfo.IPv6, cidr)
				}
			}
		}
	}

	// Save last interface info
	if currentInfo != nil {
		currentInfo.Type = determineInterfaceType(flags, media, groups)
		currentInfo.Status = determineInterfaceStatus(flags)
		result = append(result, *currentInfo)
	}

	return result
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
		return "up"
	}
	return "down"
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
