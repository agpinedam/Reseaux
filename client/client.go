package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"reseaux/rip"
	"reseaux/router"
	"reseaux/table"
)

func main() {
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"), // Dirección del servidor
		Port: rip.RIPPort,
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Leer la configuración del router desde el archivo YAML
	routerConfigPath := "../data/routeur-r2.yaml"
	r, err := router.NewRouterFromFile(routerConfigPath)
	if err != nil {
		fmt.Printf("Error loading router configuration: %v\n", err)
		return
	}

	// Construir la tabla de enrutamiento desde la configuración del router
	routeTable := table.NewRouteTableFromRouter(r)

	for {
		// Construir y enviar el mensaje RIP con la tabla de enrutamiento al servidor
		var msg rip.RIPMessage
		msg.Command = 1 // Define el comando RIP como "response"
		msg.Version = 2 // Define la versión del protocolo RIP

		// Agregar cada entrada de la tabla de enrutamiento a la estructura del mensaje RIP
		for _, iface := range routeTable.Routes {
			entry := rip.RIPEntity{
				AddressFamilyIdentifier: 2, // IPv4
				IPAddress:               iface.IP,
				SubnetMask:              net.IP(iface.Mask),
				NextHop:                 net.ParseIP("10.1.1.3"), // Dirección del servidor
				Metric:                  uint32(iface.Metric),    // Métrica inicial
			}
			msg.Entries = append(msg.Entries, entry)
		}

		data, err := msg.MarshalBinary()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		fmt.Println("Sent RIP table to server")

		// Esperar un momento antes de enviar la tabla nuevamente
		time.Sleep(time.Second * 30)
	}
}
