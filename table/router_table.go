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
		iface.Metric = 1 // Initialiser la métrique à 1
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
			Device: entry.NextHop.String(), // Ou ajuster selon l'appareil approprié
			IP:     entry.IPAddress,
			Mask:   mask,
			Metric: int(entry.Metric), // Utiliser la métrique reçue dans le message RIP
		}
	}

	return routeTable
}

func (rt *RouteTable) MergeRouteTable(newTable *RouteTable) *RouteTable {
	mergedTable := &RouteTable{
		Routes: make([]router.Interface, len(rt.Routes)),
	}

	// Copier les routes existantes
	copy(mergedTable.Routes, rt.Routes)

	// Pour chaque route dans la nouvelle table, si elle existe déjà dans la table fusionnée, l'ignorer.
	// Si elle n'existe pas, ajouter simplement la nouvelle route à la table fusionnée.
	for _, newRoute := range newTable.Routes {
		var found bool
		for _, existingRoute := range mergedTable.Routes {
			if existingRoute.Device == newRoute.Device && existingRoute.IP.Equal(newRoute.IP) && existingRoute.Mask.String() == newRoute.Mask.String() {
				found = true
				break
			}
		}
		if !found {
			newRoute.Metric++ // Incrémenter la métrique de 1 lors de l'ajout d'une nouvelle route
			mergedTable.Routes = append(mergedTable.Routes, newRoute)
		}
	}

	return mergedTable
}

func (rt *RouteTable) RecalculateMetrics() {
	// Incrémenter les métriques de toutes les routes de 1
	for i := range rt.Routes {
		rt.Routes[i].Metric++
	}
}
