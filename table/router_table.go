package table

import (
	"net"
	"reseaux/rip"
	"reseaux/router"
)

type RouteTable struct {
	Routes []router.Interface
}

func NewRouteTableFromRouter(r *router.Router) *RouteTable {
	routeTable := &RouteTable{
		Routes: make([]router.Interface, len(r.Interfaces)),
	}

	copy(routeTable.Routes, r.Interfaces)

	return routeTable
}

func BuildRouteTableFromRIPMessage(msg *rip.RIPMessage) *RouteTable {
	routeTable := &RouteTable{
		Routes: make([]router.Interface, len(msg.Entries)),
	}

	for i, entry := range msg.Entries {
		mask := net.IPMask(entry.SubnetMask)

		routeTable.Routes[i] = router.Interface{
			Device: entry.NextHop.String(), // O ajusta según el dispositivo adecuado
			IP:     entry.IPAddress,
			Mask:   mask,
			Metric: int(entry.Metric), // Convierte uint32 a int
		}
	}

	return routeTable
}

func (rt *RouteTable) MergeRouteTable(newTable *RouteTable) *RouteTable {
	mergedTable := &RouteTable{
		Routes: append(rt.Routes, newTable.Routes...),
	}

	return mergedTable
}

func (rt *RouteTable) RecalculateMetrics() {
	// Incrementar las métricas de todas las rutas en 1
	for i := range rt.Routes {
		rt.Routes[i].Metric++
	}
}
