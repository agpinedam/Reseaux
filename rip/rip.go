package rip

import (
	"encoding/binary"
	"net"
)

const (
	RIPPort = 520
)

type RIPMessage struct {
	Command uint8
	Version uint8
	Zero    uint16
	Entries []RIPEntity
}

type RIPEntity struct {
	AddressFamilyIdentifier uint16
	RouteTag                uint16
	IPAddress               net.IP
	SubnetMask              net.IP
	NextHop                 net.IP
	Metric                  uint32
}

func (msg *RIPMessage) MarshalBinary() ([]byte, error) {
	buffer := make([]byte, 4) // Encabezado RIP tiene 4 bytes
	buffer[0] = msg.Command
	buffer[1] = msg.Version
	binary.BigEndian.PutUint16(buffer[2:], msg.Zero)

	for _, entry := range msg.Entries {
		entryBytes := make([]byte, 20) // Cada entrada RIP tiene 20 bytes
		binary.BigEndian.PutUint16(entryBytes[:], entry.AddressFamilyIdentifier)
		binary.BigEndian.PutUint16(entryBytes[2:], entry.RouteTag)
		copy(entryBytes[4:], entry.IPAddress.To4())
		copy(entryBytes[8:], entry.SubnetMask.To4())
		copy(entryBytes[12:], entry.NextHop.To4())
		binary.BigEndian.PutUint32(entryBytes[16:], entry.Metric)
		buffer = append(buffer, entryBytes...)
	}

	return buffer, nil
}

func (msg *RIPMessage) UnmarshalBinary(data []byte) error {
	msg.Command = data[0]
	msg.Version = data[1]
	msg.Zero = binary.BigEndian.Uint16(data[2:4])

	data = data[4:]
	for len(data) >= 20 {
		var entry RIPEntity
		entry.AddressFamilyIdentifier = binary.BigEndian.Uint16(data[:2])
		entry.RouteTag = binary.BigEndian.Uint16(data[2:4])
		entry.IPAddress = net.IP(data[4:8])
		entry.SubnetMask = net.IP(data[8:12])
		entry.NextHop = net.IP(data[12:16])
		entry.Metric = binary.BigEndian.Uint32(data[16:20])
		msg.Entries = append(msg.Entries, entry)
		data = data[20:]
	}

	return nil
}
