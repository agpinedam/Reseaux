package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Escuchar en el puerto 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("Servidor escuchando en el puerto 8080")

	for {
		// Aceptar una conexión
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Manejar la conexión en una nueva gorutina
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Cliente conectado")
	// Enviar un mensaje al cliente
	fmt.Fprintf(conn, "Hola desde el servidor\n")
	// Leer datos del cliente
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Recibido del cliente:", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer del cliente:", err)
	}
	fmt.Println("Cliente desconectado")
}
