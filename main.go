// En el archivo main.go
package main

import (
	"fmt"
	"log"

	"reseaux/router"
)

func main() {
	// Leer el archivo YAML del router
	routerFilePath := "data/routeur-r1.yaml"
	r, err := router.NewRouterFromFile(routerFilePath)
	if err != nil {
		log.Fatalf("Error al leer el archivo del router: %v", err)
	}

	// Construir la tabla de rutas
	routeTable := router.BuildRouteTable(r)

	// Imprimir la tabla de rutas
	fmt.Println("Tabla de rutas:")
	for _, route := range routeTable.Routes {
		fmt.Printf("Dispositivo: %s, IP: %s, Máscara: %s\n", route.Device, route.IP, route.Mask)
	}
}
