package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run server.go <port>")
		os.Exit(1)
	}

	port := os.Args[1]
	runServer(port)
}

func runServer(port string) {
	addr := ":" + port
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("Servidor escuchando en el puerto %s\n", port)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error al leer del cliente:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Recibido del cliente: %s\n", message)

		response := "Respuesta desde el servidor en " + port
		_, err = conn.WriteTo([]byte(response), clientAddr)
		if err != nil {
			fmt.Println("Error al enviar respuesta al cliente:", err)
		}
	}
}
