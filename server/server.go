package main

import (
	"fmt"
	"net"

	"reseaux/rip"
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

	buffer := make([]byte, 1024)
	receivedTable := false

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		if !receivedTable {
			var msg rip.RIPMessage
			if err := msg.UnmarshalBinary(buffer[:n]); err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}

			fmt.Printf("Received RIP table from client %s: %+v\n", clientAddr, msg)
			receivedTable = true
		}
	}
}
