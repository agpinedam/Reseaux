package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func startTestServer(t *testing.T, port string) (net.PacketConn, int) {
	conn, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		t.Fatalf("No se pudo iniciar el servidor de prueba: %v", err)
	}
	actualPort := conn.LocalAddr().(*net.UDPAddr).Port

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				return
			}

			message := string(buffer[:n])
			fmt.Printf("Servidor (%s) recibió del cliente: %s\n", port, message)

			response := "Respuesta desde el servidor en " + port
			_, err = conn.WriteTo([]byte(response), addr)
			if err != nil {
				return
			}
		}
	}()

	return conn, actualPort
}

func sendTestMessage(t *testing.T, serverAddr, message string) {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		t.Fatalf("No se pudo conectar al servidor: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Error al enviar mensaje al servidor: %v", err)
	}

	reply := make([]byte, 1024)
	n, err := bufio.NewReader(conn).Read(reply)
	if err != nil {
		t.Fatalf("Error al leer la respuesta del servidor: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", string(reply[:n]))
}

func TestClient(t *testing.T) {
	serverPort := "8082"
	conn, _ := startTestServer(t, serverPort)
	defer conn.Close()

	time.Sleep(1 * time.Second) // Espera para asegurar que el servidor esté escuchando

	sendTestMessage(t, "localhost:"+serverPort, "Mensaje de prueba")
}
