package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"reseaux/utils" // Ajusta esta ruta según tu estructura de proyecto
)

func main() {
	port := "8080"
	startClient(port)
}

func startClient(port string) {
	time.Sleep(2 * time.Second)

	conn, err := net.Dial("udp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
		return
	}
	defer conn.Close()

	files := []string{
		"../data/routeur-client.yaml",
		"../data/routeur-r1.yaml",
		"../data/routeur-r2.yaml",
		"../data/routeur-r3.yaml",
		"../data/routeur-r4.yaml",
		"../data/routeur-r5.yaml",
		"../data/routeur-r6.yaml",
		"../data/routeur-serveur.yaml",
	}

	ripMessage, err := utils.GenerateRIPMessage(files)
	if err != nil {
		log.Fatalf("Error generating RIP message: %v", err)
	}

	_, err = conn.Write(ripMessage)
	if err != nil {
		log.Fatalf("Error writing: %v", err)
		return
	}

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		log.Fatalf("Error reading: %v", err)
		return
	}

	fmt.Println("Server says:", string(reply[:n]))
}
