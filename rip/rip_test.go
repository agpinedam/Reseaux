package rip_test

import (
	"net"
	"reseaux/rip"
	"testing"
)

func TestMarshalUnmarshalRIPMessage(t *testing.T) {
	// Création d'un message RIP original pour tester la sérialisation et la désérialisation
	originalMsg := &rip.RIPMessage{
		Command: 1,
		Version: 2,
		Entries: []rip.RIPEntity{
			{
				AddressFamilyIdentifier: 2,
				IPAddress:               net.ParseIP("192.168.1.1"),
				SubnetMask:              net.ParseIP("255.255.255.0").To4(), // Utilisation de net.ParseIP pour obtenir le masque de sous-réseau en tant que net.IP
				NextHop:                 net.ParseIP("192.168.1.254"),
				Metric:                  1,
			},
		},
	}

	// Sérialisation du message RIP original
	data, err := originalMsg.MarshalBinary()
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du message RIP : %v", err)
	}

	// Désérialisation du message RIP sérialisé
	var parsedMsg rip.RIPMessage
	err = parsedMsg.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("Erreur lors de la désérialisation du message RIP : %v", err)
	}

	// Vérification si le message désérialisé correspond au message original
	if parsedMsg.Command != originalMsg.Command || parsedMsg.Version != originalMsg.Version {
		t.Errorf("Le message désérialisé ne correspond pas à l'original")
	}

	// Vérification du nombre d'entrées dans le message désérialisé
	if len(parsedMsg.Entries) != len(originalMsg.Entries) {
		t.Errorf("Le nombre d'entrées désérialisées ne correspond pas à l'original")
	}

	// Comparaison de chaque entrée désérialisée avec l'entrée originale
	for i, entry := range parsedMsg.Entries {
		originalEntry := originalMsg.Entries[i]
		if !entry.IPAddress.Equal(originalEntry.IPAddress) ||
			!entry.NextHop.Equal(originalEntry.NextHop) ||
			entry.Metric != originalEntry.Metric ||
			entry.AddressFamilyIdentifier != originalEntry.AddressFamilyIdentifier ||
			entry.RouteTag != originalEntry.RouteTag ||
			entry.SubnetMask.String() != originalEntry.SubnetMask.String() {
			t.Errorf("L'entrée désérialisée ne correspond pas à l'originale")
		}
	}

	// Message de réussite
	t.Log("Sérialisation et désérialisation du message RIP réussies")
}
