package main

import (
	"fmt"
	"net"
)

func main() {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"), // Dirección de broadcast para la red local
		Port: 8080,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server listening on 10.1.1.3:8080")

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Printf("Received '%s' from %s\n", string(buffer[:n]), clientAddr)
		// No es necesario responder al cliente en un broadcast
	}
}
