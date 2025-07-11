// Package netstat provides parsing utilities for netstat output.
package netstat

import (
	"bufio"
	"fmt"
	"runtime"
	"strings"
)

// Route represents a parsed route entry from netstat output.
type Route struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Genmask     string `json:"genmask,omitempty"`
	Flags       string `json:"flags,omitempty"`
	Metric      string `json:"metric,omitempty"`
	Interface   string `json:"interface"`
}

const minRouteFields = 4

// SystemOS returns the current system's OS name in lowercase.
func SystemOS() string {
	return strings.ToLower(runtime.GOOS)
}

// ParseNetstat parses the output of 'netstat -rn' on FreeBSD and returns a slice of Route.
func ParseNetstat(output string) ([]Route, error) {
	var routes []Route
	var headerIdx = map[string]int{}
	var headerFound bool

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		// Find the header line
		if !headerFound && (strings.HasPrefix(line, "Destination") || strings.HasPrefix(line, "destination")) {
			headerFound = true
			for i, f := range fields {
				headerIdx[strings.ToLower(f)] = i
			}
			continue
		}
		if !headerFound {
			continue
		}
		// Skip lines that are not data
		if len(fields) < minRouteFields {
			continue
		}
		// Parse fields based on header
		var r Route
		if idx, ok := headerIdx["destination"]; ok && idx < len(fields) {
			r.Destination = fields[idx]
		}
		if idx, ok := headerIdx["gateway"]; ok && idx < len(fields) {
			r.Gateway = fields[idx]
		}
		if idx, ok := headerIdx["flags"]; ok && idx < len(fields) {
			r.Flags = fields[idx]
		}
		if idx, ok := headerIdx["netif"]; ok && idx < len(fields) {
			r.Interface = fields[idx]
		} else if idx, ok := headerIdx["iface"]; ok && idx < len(fields) {
			r.Interface = fields[idx]
		}
		if idx, ok := headerIdx["metric"]; ok && idx < len(fields) {
			r.Metric = fields[idx]
		}
		// FreeBSD does not have Genmask, so skip
		// Only add if Destination and Gateway are present
		if r.Destination != "" && r.Gateway != "" {
			routes = append(routes, r)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	if !headerFound {
		return nil, fmt.Errorf("no netstat header found")
	}
	return routes, nil
}
