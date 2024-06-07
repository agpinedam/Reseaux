package main

import (
	"fmt"
	"net"
	"reseaux/rip"
	"reseaux/router"
	"reseaux/table"
	"sync"
)

var (
	routeTableLock sync.Mutex
	receivedTables = make(map[string]struct{})
)

func main() {
	broadcastAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"), // Adresse de diffusion pour le réseau local
		Port: rip.RIPPort,
	}

	conn, err := net.ListenUDP("udp", broadcastAddr)
	if err != nil {
		fmt.Printf("Erreur : %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Serveur UDP à l'écoute sur 10.1.1.3:520")

	// Lire la configuration du routeur depuis le fichier YAML
	cheminConfigRouteur := "../data/routeur-r1.yaml"
	r, err := router.NewRouterFromFile(cheminConfigRouteur)
	if err != nil {
		fmt.Printf("Erreur lors du chargement de la configuration du routeur : %v\n", err)
		return
	}

	// Construire la table de routage à partir de la configuration du routeur
	routeTable := table.NewRouteTableFromRouter(r)
	fmt.Printf("Table de routage du serveur : %+v\n", routeTable)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Erreur : %s\n", err)
			continue
		}

		var msg rip.RIPMessage
		if err := msg.UnmarshalBinary(buffer[:n]); err != nil {
			fmt.Printf("Erreur : %s\n", err)
			continue
		}

		tableID := fmt.Sprintf("%s:%d", clientAddr.IP, clientAddr.Port)

		// Vérifier si cette table a déjà été reçue
		routeTableLock.Lock()
		if _, ok := receivedTables[tableID]; ok {
			routeTableLock.Unlock()
			continue
		}
		receivedTables[tableID] = struct{}{}
		routeTableLock.Unlock()

		fmt.Printf("Table RIP reçue du client %s : %+v\n", clientAddr, msg)

		// Construire la table de routage à partir du message RIP reçu
		receivedRouteTable := table.BuildRouteTableFromRIPMessage(&msg)

		// Fusionner la table de routage reçue avec la table du serveur
		routeTableLock.Lock()
		routeTable = routeTable.MergeRouteTable(receivedRouteTable)
		routeTableLock.Unlock()

		fmt.Printf("Table de routage fusionnée : %+v\n", routeTable)
	}
}
