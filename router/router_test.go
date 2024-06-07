package router

import (
	"net"
	"testing"
)

func TestNewRouterFromFile(t *testing.T) {
	routerConfigPath := "../data/routeur-r1.yaml"
	r, err := NewRouterFromFile(routerConfigPath)
	if err != nil {
		t.Fatalf("Erreur lors du chargement de la configuration du routeur: %v", err)
	}

	if len(r.Interfaces) == 0 {
		t.Error("La configuration du routeur ne contient pas d'interfaces")
	}

	t.Log("Configuration du routeur chargée correctement")
}

func TestInterfaceParsing(t *testing.T) {
	rawIface := rawInterface{
		Device: "eth0",
		IP:     "192.168.1.1",
		Mask:   24,
	}

	ip := net.ParseIP(rawIface.IP)
	if ip == nil {
		t.Errorf("Adresse IP invalide: %s", rawIface.IP)
	}

	mask := net.CIDRMask(rawIface.Mask, 32)
	if len(mask) == 0 {
		t.Errorf("Masque de sous-réseau invalide: %d", rawIface.Mask)
	}

	t.Log("Analyse de l'interface réussie")
}
