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

	for i, iface := range r.Interfaces {
		iface.Metric = 1 // Inicializar la métrica en 1
		routeTable.Routes[i] = iface
	}

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
			Metric: int(entry.Metric), // Utiliza la métrica recibida en el mensaje RIP
		}
	}

	return routeTable
}

func (rt *RouteTable) MergeRouteTable(newTable *RouteTable) *RouteTable {
	mergedTable := &RouteTable{
		Routes: make([]router.Interface, len(rt.Routes)),
	}

	// Copiar las rutas existentes
	copy(mergedTable.Routes, rt.Routes)

	// Para cada ruta en la nueva tabla, si ya existe en la tabla fusionada, ignórela.
	// Si no existe, simplemente agrega la nueva ruta a la tabla fusionada.
	for _, newRoute := range newTable.Routes {
		var found bool
		for _, existingRoute := range mergedTable.Routes {
			if existingRoute.Device == newRoute.Device && existingRoute.IP.Equal(newRoute.IP) && existingRoute.Mask.String() == newRoute.Mask.String() {
				found = true
				break
			}
		}
		if !found {
			newRoute.Metric++ // Incrementar la métrica en 1 al añadir una nueva ruta
			mergedTable.Routes = append(mergedTable.Routes, newRoute)
		}
	}

	return mergedTable
}

func (rt *RouteTable) RecalculateMetrics() {
	// Incrementar las métricas de todas las rutas en 1
	for i := range rt.Routes {
		rt.Routes[i].Metric++
	}
}
