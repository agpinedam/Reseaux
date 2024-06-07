package table

import (
	"net"
	"reseaux/rip"
	"reseaux/router"
	"testing"
)

func TestNewRouteTableFromRouter(t *testing.T) {
	r := &router.Router{
		Interfaces: []router.Interface{
			{Device: "eth0", IP: net.ParseIP("192.168.1.1"), Mask: net.CIDRMask(24, 32)},
		},
	}

	routeTable := NewRouteTableFromRouter(r)
	if len(routeTable.Routes) == 0 {
		t.Error("La table de routage n'a pas été correctement construite à partir de la configuration du routeur")
	}

	t.Log("Table de routage construite correctement")
}

func TestMergeRouteTable(t *testing.T) {
	rt1 := &RouteTable{
		Routes: []router.Interface{
			{Device: "eth0", IP: net.ParseIP("192.168.1.1"), Mask: net.CIDRMask(24, 32), Metric: 1},
		},
	}

	rt2 := &RouteTable{
		Routes: []router.Interface{
			{Device: "eth1", IP: net.ParseIP("192.168.2.1"), Mask: net.CIDRMask(24, 32), Metric: 1},
		},
	}

	mergedTable := rt1.MergeRouteTable(rt2)
	if len(mergedTable.Routes) != 2 {
		t.Error("La table de routage fusionnée ne contient pas les routes attendues")
	}

	t.Log("Table de routage fusionnée correctement")
}

func TestBuildRouteTableFromRIPMessage(t *testing.T) {
	msg := &rip.RIPMessage{
		Command: 1,
		Version: 2,
		Entries: []rip.RIPEntity{
			{
				AddressFamilyIdentifier: 2,
				IPAddress:               net.ParseIP("192.168.3.1"),
				SubnetMask:              net.IPv4(255, 255, 255, 0),
				NextHop:                 net.ParseIP("192.168.3.254"),
				Metric:                  1,
			},
		},
	}

	routeTable := BuildRouteTableFromRIPMessage(msg)
	if len(routeTable.Routes) == 0 {
		t.Error("La table de routage n'a pas été correctement construite à partir du message RIP")
	}

	t.Log("Table de routage construite correctement à partir du message RIP")
}
