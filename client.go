package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Uso: go run client.go <local-port> <remote-address> <message>")
		os.Exit(1)
	}

	localPort := os.Args[1]
	remoteAddress := os.Args[2]
	message := os.Args[3]

	go runServer(localPort)
	time.Sleep(100 * time.Millisecond) // Espera para asegurar que el servidor esté escuchando

	sendMessage(remoteAddress, message)
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

func sendMessage(serverAddr, message string) {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		fmt.Println("Error al conectar con el servidor:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error al escribir:", err)
		return
	}

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return
	}

	fmt.Printf("El servidor responde: %s\n", string(reply[:n]))
}
