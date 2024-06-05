package table

import "reseaux/router"

type RouteTable struct {
	Routes []router.Interface
}

func BuildRouteTable(r *router.Router) *RouteTable {
	routeTable := &RouteTable{
		Routes: r.Interfaces,
	}

	return routeTable
}
