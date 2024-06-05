package main

import (
	"fmt"
	"log"
	"net"

	"reseaux/router"
	"reseaux/table"
)

func main() {
	routerConfigPath := "data/routeur-r1.yaml"
	r, err := router.NewRouterFromFile(routerConfigPath)
	if err != nil {
		log.Fatalf("Failed to load router configuration: %v", err)
	}

	routeTable := table.BuildRouteTable(r)

	fmt.Println("Router Interfaces:")
	for _, iface := range routeTable.Routes {
		fmt.Printf("Device: %s, IP: %s, Mask: %s\n", iface.Device, iface.IP, net.IP(iface.Mask))
	}
}
