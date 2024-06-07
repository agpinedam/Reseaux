package main

import (
	"net"
	"reseaux/rip"
	"reseaux/router"
	"reseaux/table"
	"testing"
)

func TestServerInitialization(t *testing.T) {
	// Configurer le serveur UDP
	broadcastAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"),
		Port: rip.RIPPort,
	}

	conn, err := net.ListenUDP("udp", broadcastAddr)
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation du serveur UDP: %v", err)
	}
	defer conn.Close()

	t.Log("Serveur UDP démarré correctement sur 10.1.1.3:520")
}

func TestRouterConfigLoading(t *testing.T) {
	routerConfigPath := "../data/routeur-r1.yaml"
	r, err := router.NewRouterFromFile(routerConfigPath)
	if err != nil {
		t.Fatalf("Erreur lors du chargement de la configuration du routeur: %v", err)
	}

	if len(r.Interfaces) == 0 {
		t.Error("La configuration du routeur ne contient pas d'interfaces")
	}

	t.Log("Configuration du routeur chargée correctement")
}

func TestRouteTableConstruction(t *testing.T) {
	r := &router.Router{
		Interfaces: []router.Interface{
			{Device: "eth0", IP: net.ParseIP("192.168.1.1"), Mask: net.CIDRMask(24, 32)},
		},
	}

	routeTable := table.NewRouteTableFromRouter(r)
	if len(routeTable.Routes) == 0 {
		t.Error("La table de routage n'a pas été correctement construite à partir de la configuration du routeur")
	}

	t.Log("Table de routage construite correctement")
}
