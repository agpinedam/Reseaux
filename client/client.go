package main

import (
	"fmt"
	"net"
	"os"

	"reseaux/rip"
)

func main() {
	broadcastAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"), // Dirección de broadcast para la red local
		Port: rip.RIPPort,
	}

	conn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Construir y enviar la tabla de enrutamiento RIP al servidor
	var msg rip.RIPMessage
	// Aquí debes construir la tabla de enrutamiento RIP y agregarla al mensaje RIP
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
}
