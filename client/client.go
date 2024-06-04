package main

import (
	"fmt"
	"net"
	"time"
)

// Función para iniciar el cliente UDP
func startClient(port string) {
	time.Sleep(2 * time.Second) // Espera a que el servidor se inicie completamente

	conn, err := net.Dial("udp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	message := "Hello, server!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	fmt.Println("Server says:", string(reply))
}

func main() {
	port := "8080"
	startClient(port)
}
