package main

import (
	"fmt"
	"net"
)

func main() {
	sendMessage("localhost:8080", "Hello, server!")
}

func sendMessage(serverAddr, message string) {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}

	reply := make([]byte, 4096)
	_, err = conn.Read(reply)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	fmt.Println("Server says:", string(reply))
}
