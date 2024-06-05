package main

import (
	"fmt"
	"net"
)

func main() {
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("10.1.1.1"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server listening on 10.1.1.1:8080")

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Printf("Received '%s' from %s\n", string(buffer[:n]), clientAddr)
		conn.WriteToUDP([]byte("Hello from server"), clientAddr)
	}
}
