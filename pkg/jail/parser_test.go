package jail

import (
	"reflect"
	"testing"
)

func TestParseJailList(t *testing.T) {
	t.Run("parses typical output with IPv4 and IPv6", func(t *testing.T) {
		output := `name       state   ip4.addr      ip6.addr         path
jail1      running  192.168.1.10  -               /jails/jail1
jail2      stopped  10.0.0.2      2001:db8::2     /jails/jail2
jail3      running  -             fe80::1         /jails/jail3`
		want := []Info{
			{Name: "jail1", Status: "running", IPv4: "192.168.1.10", IPv6: "-", Path: "/jails/jail1"},
			{Name: "jail2", Status: "stopped", IPv4: "10.0.0.2", IPv6: "2001:db8::2", Path: "/jails/jail2"},
			{Name: "jail3", Status: "running", IPv4: "-", IPv6: "fe80::1", Path: "/jails/jail3"},
		}
		got, err := ParseJailList(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("empty output returns empty slice", func(t *testing.T) {
		output := ""
		got, err := ParseJailList(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %+v", got)
		}
	})

	t.Run("malformed lines are skipped", func(t *testing.T) {
		output := `name state ip4.addr ip6.addr path
jail1 running 192.168.1.10 - /jails/jail1
malformed line here
jail2 stopped 10.0.0.2 2001:db8::2 /jails/jail2`
		want := []Info{
			{Name: "jail1", Status: "running", IPv4: "192.168.1.10", IPv6: "-", Path: "/jails/jail1"},
			{Name: "jail2", Status: "stopped", IPv4: "10.0.0.2", IPv6: "2001:db8::2", Path: "/jails/jail2"},
		}
		got, err := ParseJailList(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
