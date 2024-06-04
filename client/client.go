package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

// Función para iniciar el cliente UDP
func startClient(port string) {
	time.Sleep(2 * time.Second) // Espera a que el servidor se inicie completamente

	conn, err := net.Dial("udp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	// Leer el archivo RIP
	ripMessage, err := ioutil.ReadFile("../rip_message.bin")
	if err != nil {
		fmt.Println("Error reading RIP message:", err.Error())
		return
	}

	// Enviar el mensaje RIP
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

func main() {
	port := "8080"
	startClient(port)
}
