package rip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

type RIPPacket struct {
	Command    uint8
	Version    uint8
	Zero       uint16
	RIPEntries []RIPEntry
}

type RIPEntry struct {
	AFI        uint16
	RouteTag   uint16
	IPAddress  uint32
	SubnetMask uint32
	NextHop    uint32
	Metric     uint32
}

func CreateRIPRequest() *RIPPacket {
	return &RIPPacket{
		Command: 1, // Request
		Version: 2,
		RIPEntries: []RIPEntry{
			{
				AFI:    2,
				Metric: 16,
			},
		},
	}
}

func (r *RIPPacket) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SendRIPRequest(conn *net.UDPConn) error {
	ripPacket := CreateRIPRequest()
	data, err := ripPacket.Serialize()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func ListenForRIPResponses(conn *net.UDPConn) {
	buf := make([]byte, 512)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s\n", n, addr)
		// Procesar la respuesta RIP aquí
	}
}

func StartRIPProtocol(iface string, port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(iface),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go ListenForRIPResponses(conn)

	for {
		err = SendRIPRequest(conn)
		if err != nil {
			fmt.Println("Error sending RIP request:", err)
		}
		time.Sleep(30 * time.Second)
	}
}
