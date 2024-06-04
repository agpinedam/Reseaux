package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	port := "8080" // Asegúrate de que este puerto coincida con el utilizado en el cliente y el servidor

	// Ejecutar el servidor en una goroutine
	go func() {
		cmd := exec.Command("go", "run", "server/server.go", port)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error starting server:", err)
			os.Exit(1)
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println("Server exited with error:", err)
		}
	}()

	// Esperar a que el servidor se inicie completamente
	time.Sleep(2 * time.Second)

	// Ejecutar el cliente
    go startClient(8080,"/");
	//cmd := exec.Command("go", "run", "client/client.go", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running client:", err)
		os.Exit(1)
	}
}
