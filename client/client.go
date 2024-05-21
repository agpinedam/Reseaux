package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Conectarse al servidor en el puerto 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Conectado al servidor")

	// Leer mensaje del servidor
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer del servidor:", err)
		return
	}
	fmt.Print("Mensaje del servidor: ", message)

	// Enviar mensaje al servidor
	fmt.Fprintf(conn, "Hola desde el cliente\n")

	// Leer respuestas adicionales del servidor
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Recibido del servidor:", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer del servidor:", err)
	}
}
