// En el archivo router/router_table.go
package router

type RouteTable struct {
	Routes []Interface
}

func BuildRouteTable(r *Router) *RouteTable {
	routeTable := &RouteTable{}

	for _, iface := range r.Interfaces {
		// Build routes for each interface
		routeTable.Routes = append(routeTable.Routes, iface)
	}

	return routeTable
}
