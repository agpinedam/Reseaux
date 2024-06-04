package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
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

			message := string(buffer[:n])
			fmt.Println("Recibido del cliente:", message)

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

	message := "Mensaje de prueba"
	_, err = clientConn.Write([]byte(message))
	if err != err {
		t.Fatalf("Error al enviar mensaje al servidor: %v", err)
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
