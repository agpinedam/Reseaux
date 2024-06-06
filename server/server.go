package main

import (
	"fmt"
	"net"

	"reseaux/rip"
	"reseaux/router"
	"reseaux/table"
)

func main() {
	broadcastAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"), // Dirección de broadcast para la red local
		Port: rip.RIPPort,
	}

	conn, err := net.ListenUDP("udp", broadcastAddr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server listening on 10.1.1.3:520")

	// Leer la configuración del router desde el archivo YAML
	routerConfigPath := "../data/routeur-r1.yaml"
	r, err := router.NewRouterFromFile(routerConfigPath)
	if err != nil {
		fmt.Printf("Error loading router configuration: %v\n", err)
		return
	}

	// Construir la tabla de enrutamiento desde la configuración del router
	routeTable := table.NewRouteTableFromRouter(r)
	fmt.Printf("Server routing table: %+v\n", routeTable)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		var msg rip.RIPMessage
		if err := msg.UnmarshalBinary(buffer[:n]); err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		fmt.Printf("Received RIP table from client %s: %+v\n", clientAddr, msg)

		// Crear una nueva tabla de enrutamiento desde el mensaje RIP recibido
		newRouteTable := table.BuildRouteTableFromRIPMessage(&msg)

		// Fusionar la nueva tabla con la tabla existente
		mergedRouteTable := routeTable.MergeRouteTable(newRouteTable)

		// Recalcular las métricas en la tabla fusionada
		mergedRouteTable.RecalculateMetrics()

		// Actualizar la tabla de enrutamiento del servidor con la tabla fusionada
		routeTable = mergedRouteTable

		fmt.Printf("Merged routing table: %+v\n", routeTable)
	}
}
