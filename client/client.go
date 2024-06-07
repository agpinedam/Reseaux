package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"reseaux/rip"
	"reseaux/router"
	"reseaux/table"
)

func main() {
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"), // Adresse du serveur
		Port: rip.RIPPort,
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Printf("Erreur : %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Lire la configuration du routeur depuis le fichier YAML
	cheminConfigRouteur := "../data/routeur-r2.yaml"
	r, err := router.NewRouterFromFile(cheminConfigRouteur)
	if err != nil {
		fmt.Printf("Erreur lors du chargement de la configuration du routeur : %v\n", err)
		return
	}

	// Construire la table de routage à partir de la configuration du routeur
	routeTable := table.NewRouteTableFromRouter(r)

	for {
		// Construire et envoyer le message RIP avec la table de routage au serveur
		var msg rip.RIPMessage
		msg.Command = 1 // Définir la commande RIP comme "réponse"
		msg.Version = 2 // Définir la version du protocole RIP

		// Ajouter chaque entrée de la table de routage à la structure du message RIP
		for _, iface := range routeTable.Routes {
			entry := rip.RIPEntity{
				AddressFamilyIdentifier: 2, // IPv4
				IPAddress:               iface.IP.Mask(iface.Mask),
				SubnetMask:              net.IP(iface.Mask),
				NextHop:                 net.ParseIP("10.1.1.3"), // Adresse du serveur
				Metric:                  uint32(iface.Metric),    // Métrique initiale
			}
			msg.Entries = append(msg.Entries, entry)
		}

		data, err := msg.MarshalBinary()
		if err != nil {
			fmt.Printf("Erreur : %s\n", err)
			return
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Printf("Erreur : %s\n", err)
			return
		}

		fmt.Println("Table RIP envoyée au serveur")

		// Attendre un moment avant d'envoyer la table à nouveau
		time.Sleep(time.Second * 30)
	}

}
