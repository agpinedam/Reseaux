package router

type RouteTable struct {
	Routes []Interface
}

func BuildRouteTable(r *Router) *RouteTable {
	routeTable := &RouteTable{
		Routes: r.Interfaces,
	}

	return routeTable
}
