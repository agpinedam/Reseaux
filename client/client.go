package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func main() {
	port := "8080"
	startClient(port)
}

func startClient(port string) {
	time.Sleep(2 * time.Second)

	conn, err := net.Dial("udp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	ripMessage, err := ioutil.ReadFile("rip_message.bin")
	if err != nil {
		fmt.Println("Error reading RIP message:", err.Error())
		return
	}

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
