package netstat

import (
	"reflect"
	"testing"
)

func TestParseNetstat_FreeBSD(t *testing.T) {
	sample := `Routing tables

Internet:
Destination        Gateway            Flags     Netif Expire
10.0.0.0/24        10.0.0.1           UGS      em0
10.0.0.1           link#1             UHS      lo0
127.0.0.1          link#2             UH       lo0
192.168.1.0/24     192.168.1.1        UGS      em1
192.168.1.1        link#3             UHS      lo0
`

	expected := []Route{
		{Destination: "10.0.0.0/24", Gateway: "10.0.0.1", Flags: "UGS", Interface: "em0"},
		{Destination: "10.0.0.1", Gateway: "link#1", Flags: "UHS", Interface: "lo0"},
		{Destination: "127.0.0.1", Gateway: "link#2", Flags: "UH", Interface: "lo0"},
		{Destination: "192.168.1.0/24", Gateway: "192.168.1.1", Flags: "UGS", Interface: "em1"},
		{Destination: "192.168.1.1", Gateway: "link#3", Flags: "UHS", Interface: "lo0"},
	}

	routes, err := ParseNetstat(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(routes, expected) {
		t.Errorf("parsed routes do not match expected.\nGot: %#v\nWant: %#v", routes, expected)
	}
}

func TestParseNetstat_EmptyInput(t *testing.T) {
	routes, err := ParseNetstat("")
	if err == nil {
		t.Error("expected error for missing header, got nil")
	}
	if len(routes) != 0 {
		t.Errorf("expected no routes, got: %#v", routes)
	}
}

func TestParseNetstat_MalformedLines(t *testing.T) {
	sample := `Routing tables
Internet:
Destination        Gateway            Flags     Netif Expire
MALFORMED LINE
10.0.0.0/24        10.0.0.1           UGS      em0
`
	expected := []Route{
		{Destination: "10.0.0.0/24", Gateway: "10.0.0.1", Flags: "UGS", Interface: "em0"},
	}
	routes, err := ParseNetstat(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(routes, expected) {
		t.Errorf("parsed routes do not match expected.\nGot: %#v\nWant: %#v", routes, expected)
	}
}
