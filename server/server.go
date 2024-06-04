package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Escuchar en el puerto 8080
	conn, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Servidor escuchando en el puerto 8080")

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Cliente conectado")

		// Enviar un mensaje al cliente
		_, err = conn.WriteTo([]byte("Hola desde el servidor"), addr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		message := string(buffer[:n])
		fmt.Println("Recibido del cliente:", message)
	}
}
