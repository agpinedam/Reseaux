package main

import (
	"fmt"
	"log"
	"net"

	"reseaux/router"
	"reseaux/table"
)

func calculateNetwork(ip net.IP, mask net.IPMask) net.IP {
	ip = ip.To4()
	if ip == nil {
		return nil
	}

	network := make(net.IP, len(ip))
	for i := 0; i < len(ip); i++ {
		network[i] = ip[i] & mask[i]
	}
	return network
}

func main() {
	routerConfigPath := "data/routeur-r1.yaml"
	r, err := router.NewRouterFromFile(routerConfigPath)
	if err != nil {
		log.Fatalf("Failed to load router configuration: %v", err)
	}

	routeTable := table.BuildRouteTable(r)

	fmt.Println("IP Destination | Masque de sous-réseau | Passerelle | Interface | Métrique")
	for _, iface := range routeTable.Routes {
		network := calculateNetwork(iface.IP, iface.Mask)
		fmt.Printf("%s | %s | - | %s | 1\n", network, net.IP(iface.Mask).String(), iface.IP.String())
	}
}
