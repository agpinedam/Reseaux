package main

import (
	"net"
	"reseaux/rip"
	"reseaux/router"
	"reseaux/table"
	"testing"
)

func TestClientInitialization(t *testing.T) {
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"),
		Port: rip.RIPPort,
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		t.Fatalf("Erreur lors de la connexion au serveur UDP: %v", err)
	}
	defer conn.Close()

	t.Log("Client UDP connecté correctement au serveur sur 10.1.1.3:520")
}

func TestSendRIPMessage(t *testing.T) {
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.1.1.3"),
		Port: rip.RIPPort,
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		t.Fatalf("Erreur lors de la connexion au serveur UDP: %v", err)
	}
	defer conn.Close()

	r := &router.Router{
		Interfaces: []router.Interface{
			{Device: "eth0", IP: net.ParseIP("192.168.1.1"), Mask: net.CIDRMask(24, 32), Metric: 1},
		},
	}

	routeTable := table.NewRouteTableFromRouter(r)

	var msg rip.RIPMessage
	msg.Command = 1
	msg.Version = 2

	for _, iface := range routeTable.Routes {
		entry := rip.RIPEntity{
			AddressFamilyIdentifier: 2,
			IPAddress:               iface.IP.Mask(iface.Mask),
			SubnetMask:              net.IP(iface.Mask),
			NextHop:                 net.ParseIP("10.1.1.3"),
			Metric:                  uint32(iface.Metric),
		}
		msg.Entries = append(msg.Entries, entry)
	}

	data, err := msg.MarshalBinary()
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du message RIP: %v", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		t.Fatalf("Erreur lors de l'envoi du message RIP: %v", err)
	}

	t.Log("Message RIP envoyé correctement au serveur")
}
