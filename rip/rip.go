package rip

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	// RIPPort es el puerto UDP estándar para el protocolo RIP
	RIPPort = 520
)

// RIPMessage representa un mensaje RIP.
type RIPMessage struct {
	Command byte        // Comando RIP (1 para respuesta, 2 para solicitud)
	Version byte        // Versión del protocolo RIP
	Zero    uint16      // Campo de relleno
	Entries []RIPEntity // Entradas de la tabla de enrutamiento RIP
}

// RIPEntity representa una entrada de la tabla de enrutamiento RIP.
type RIPEntity struct {
	AddressFamilyIdentifier uint16 // Identificador de la familia de direcciones (siempre 2 para IPv4)
	RouteTag                uint16 // Etiqueta de ruta (siempre 0)
	IPAddress               net.IP // Dirección IP de destino
	SubnetMask              net.IP // Máscara de subred
	NextHop                 net.IP // Próximo salto
	Metric                  uint32 // Métrica asociada a la ruta
}

// MarshalBinary codifica el mensaje RIP en formato binario.
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

// UnmarshalBinary decodifica el mensaje RIP binario en la estructura de mensaje RIP.
func (msg *RIPMessage) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("mensaje RIP demasiado corto")
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
