package rip

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	// RIPPort est le port UDP standard pour le protocole RIP
	RIPPort = 520
)

// RIPMessage représente un message RIP.
type RIPMessage struct {
	Command byte        // Commande RIP (1 pour réponse, 2 pour demande)
	Version byte        // Version du protocole RIP
	Zero    uint16      // Champ de remplissage
	Entries []RIPEntity // Entrées de la table de routage RIP
}

// RIPEntity représente une entrée de la table de routage RIP.
type RIPEntity struct {
	AddressFamilyIdentifier uint16 // Identificateur de la famille d'adresses (toujours 2 pour IPv4)
	RouteTag                uint16 // Étiquette de route (toujours 0)
	IPAddress               net.IP // Adresse IP de destination
	SubnetMask              net.IP // Masque de sous-réseau
	NextHop                 net.IP // Prochain saut
	Metric                  uint32 // Métrique associée à la route
}

// MarshalBinary encode le message RIP en format binaire.
func (msg *RIPMessage) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 4+len(msg.Entries)*20)

	buf[0] = msg.Command
	buf[1] = msg.Version
	binary.BigEndian.PutUint16(buf[2:], msg.Zero)

	offset := 4
	for _, entry := range msg.Entries {
		binary.BigEndian.PutUint16(buf[offset:], entry.AddressFamilyIdentifier)
		binary.BigEndian.PutUint16(buf[offset+2:], entry.RouteTag)
		copy(buf[offset+4:offset+8], entry.IPAddress.To4())
		copy(buf[offset+8:offset+12], entry.SubnetMask)
		copy(buf[offset+12:offset+16], entry.NextHop.To4())
		binary.BigEndian.PutUint32(buf[offset+16:], entry.Metric)
		offset += 20
	}

	return buf, nil
}

// UnmarshalBinary décode le message RIP binaire en structure de message RIP.
func (msg *RIPMessage) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("message RIP trop court")
	}

	msg.Command = data[0]
	msg.Version = data[1]
	msg.Zero = binary.BigEndian.Uint16(data[2:4])

	entriesLen := (len(data) - 4) / 20
	msg.Entries = make([]RIPEntity, entriesLen)

	offset := 4
	for i := 0; i < entriesLen; i++ {
		msg.Entries[i].AddressFamilyIdentifier = binary.BigEndian.Uint16(data[offset:])
		msg.Entries[i].RouteTag = binary.BigEndian.Uint16(data[offset+2:])
		msg.Entries[i].IPAddress = net.IP(data[offset+4 : offset+8])
		msg.Entries[i].SubnetMask = net.IP(data[offset+8 : offset+12])
		msg.Entries[i].NextHop = net.IP(data[offset+12 : offset+16])
		msg.Entries[i].Metric = binary.BigEndian.Uint32(data[offset+16 : offset+20])
		offset += 20
	}

	return nil
}
