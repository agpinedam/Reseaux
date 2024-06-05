package main

import (
	"fmt"
	"log"
	"net"

	"reseaux/utils" // Ajusta esta ruta según tu estructura de proyecto
)

func main() {
	files := []string{
		"data/routeur-client.yaml",
		"data/routeur-r1.yaml",
		"data/routeur-r2.yaml",
		"data/routeur-r3.yaml",
		"data/routeur-r4.yaml",
		"data/routeur-r5.yaml",
		"data/routeur-r6.yaml",
		"data/routeur-serveur.yaml",
	}

	ripMessage, err := utils.GenerateRIPMessage(files)
	if err != nil {
		log.Fatalf("Error al generar el mensaje RIP: %v", err)
	}

	err = sendRIPMessage(ripMessage, "localhost:8080")
	if err != nil {
		log.Fatalf("Error al enviar el mensaje RIP: %v", err)
	}

	fmt.Println("Mensaje RIP enviado correctamente")
}

func sendRIPMessage(message []byte, address string) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return fmt.Errorf("error conectando al servidor: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		return fmt.Errorf("error enviando el mensaje: %v", err)
	}
	return nil
}
