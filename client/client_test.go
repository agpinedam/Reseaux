package main

import (
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	addr := ":8082"
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Canal para enviar errores desde gorutinas secundarias
	errCh := make(chan error)

	go func() {
		buffer := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			errCh <- err
			return
		}

		message := string(buffer[:n])
		expected := "Hello, server!"
		if strings.Compare(expected, message) != 0 {
			errCh <- errors.New("Unexpected message")
			return
		}

		_, err = conn.WriteTo([]byte("Server reply"), clientAddr)
		if err != nil {
			errCh <- err
			return
		}

		// Cerramos el canal para indicar que la gorutina ha terminado
		close(errCh)
	}()

	// Espera un breve tiempo para que el servidor se inicie
	time.Sleep(100 * time.Millisecond)

	// Ejecuta la función de prueba
	sendMessage("localhost:8082", "Hello, server!")

	// Manejo de errores de gorutinas secundarias
	select {
	case err, ok := <-errCh:
		if ok && err != nil {
			t.Fatal(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout esperando respuesta del servidor")
	}
}
