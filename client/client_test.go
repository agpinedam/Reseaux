package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
	// Ajusta esta ruta según tu estructura de proyecto
)

func startTestServer(t *testing.T) (net.PacketConn, int) {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		t.Fatalf("No se pudo iniciar el servidor de prueba: %v", err)
	}
	port := conn.LocalAddr().(*net.UDPAddr).Port

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				return
			}

			fmt.Printf("Mensaje RIP recibido en hexadecimal: % X\n", buffer[:n])
			fmt.Printf("Mensaje RIP recibido en octetos: %v\n", buffer[:n])

			response := "Hola desde el servidor"
			_, err = conn.WriteTo([]byte(response), addr)
			if err != nil {
				return
			}
		}
	}()

	// Pequeña espera para asegurar que el servidor de prueba esté completamente iniciado
	time.Sleep(500 * time.Millisecond)

	return conn, port
}

func TestServer(t *testing.T) {
	conn, port := startTestServer(t)
	defer conn.Close()

	clientConn, err := net.Dial("udp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("No se pudo conectar al servidor: %v", err)
	}
	defer clientConn.Close()

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
		t.Fatalf("Error generating RIP message: %v", err)
	}

	_, err = clientConn.Write(ripMessage)
	if err != nil {
		t.Fatalf("Error al enviar mensaje RIP al servidor: %v", err)
	}

	reply := make([]byte, 1024)
	n, err := bufio.NewReader(clientConn).Read(reply)
	if err != nil {
		t.Fatalf("Error al leer la respuesta del servidor: %v", err)
	}

	expectedMessage := "Hola desde el servidor"
	if string(reply[:n]) != expectedMessage {
		t.Errorf("Mensaje esperado: %s, pero se recibió: %s", expectedMessage, string(reply[:n]))
	}
}
