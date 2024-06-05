package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	broadcastAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"), // Dirección de broadcast para la red local
		Port: 8080,
	}

	conn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	message := []byte("Hello from client")
	_, err = conn.Write(message)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Message sent to broadcast address")
}
