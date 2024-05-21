// client_test.go

package main

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"testing"
)

func TestSendMessage(t *testing.T) {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// Canal para enviar errores desde gorutinas secundarias
	errCh := make(chan error)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer conn.Close()

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}

		expected := "Hello, server!\n"
		if strings.Compare(expected, message) != 0 {
			errCh <- errors.New("Unexpected message")
			return
		}

		_, err = conn.Write([]byte("Server reply\n"))
		if err != nil {
			errCh <- err
			return
		}
	}()

	// Ejecuta la función de prueba
	sendMessage("localhost:8082", "Hello, server!")

	// Manejo de errores de gorutinas secundarias
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	default:
		// No hay errores, continuar con el test
	}
}
