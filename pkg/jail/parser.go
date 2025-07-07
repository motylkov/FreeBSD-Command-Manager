// Package jail provides utilities for parsing FreeBSD jail command output.
package jail

import (
	"bufio"
	"strings"
)

// Info represents information about a jail.
type Info struct {
	Name   string
	Status string
	IPv4   string // IPv4
	IPv6   string // IPv6
	Path   string
}

// ParseJailList parses the output of 'jail -l' and returns a slice of Info.
func ParseJailList(output string) ([]Info, error) {
	var jails []Info
	scanner := bufio.NewScanner(strings.NewReader(output))
	var headers []string
	for scanner.Scan() {
		line := scanner.Text()
		if len(headers) == 0 {
			headers = strings.Fields(line)
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != len(headers) {
			continue
		}
		jail := Info{}
		for i, h := range headers {
			switch h {
			case "name":
				jail.Name = fields[i]
			case "state":
				jail.Status = fields[i]
			case "ip4.addr":
				jail.IPv4 = fields[i]
			case "ip6.addr":
				jail.IPv6 = fields[i]
			case "path":
				jail.Path = fields[i]
			}
		}
		jails = append(jails, jail)
	}
	return jails, nil
}
