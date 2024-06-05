package main

import (
	"fmt"
	"net"
	"time"
	// Ajusta esta ruta según tu estructura de proyecto
)

func main() {
	port := "8080"
	startClient(port)
}

func startClient(port string) {
	time.Sleep(2 * time.Second)

	conn, err := net.Dial("udp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	ripMessage, err := generateRIPMessage([]string{
		"data/routeur-client.yaml",
		"data/routeur-r1.yaml",
		"data/routeur-r2.yaml",
		"data/routeur-r3.yaml",
		"data/routeur-r4.yaml",
		"data/routeur-r5.yaml",
		"data/routeur-r6.yaml",
		"data/routeur-serveur.yaml",
	})
	if err != nil {
		fmt.Println("Error generating RIP message:", err.Error())
		return
	}

	_, err = conn.Write(ripMessage)
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	fmt.Println("Server says:", string(reply[:n]))
}
