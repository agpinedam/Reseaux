package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Interfaz struct {
	Dispositivo string `yaml:"device"`
	IP          string `yaml:"ip"`
	Mascara     string `yaml:"mask"`
}

func ipToBytes(ip string) []byte {
	parsedIP := net.ParseIP(ip).To4()
	return parsedIP
}

func maskToBytes(mask string) []byte {
	var n int
	fmt.Sscanf(mask, "%d", &n)
	maskInt := ^uint32(0) << (32 - n)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, maskInt)
	return b
}

func cidrToSubnetMask(mask string) string {
	var n int
	fmt.Sscanf(mask, "%d", &n)
	maskInt := ^uint32(0) << (32 - n)
	return fmt.Sprintf("%d.%d.%d.%d", byte(maskInt>>24), byte(maskInt>>16), byte(maskInt>>8), byte(maskInt))
}

func calculateNetworkAddress(ip, mask string) string {
	ipBytes := ipToBytes(ip)
	maskBytes := maskToBytes(mask)

	network := make([]byte, 4)
	for i := 0; i < 4; i++ {
		network[i] = ipBytes[i] & maskBytes[i]
	}

	return fmt.Sprintf("%d.%d.%d.%d", network[0], network[1], network[2], network[3])
}

func printRIPFormat(interfaces []struct{ Interface Interfaz }, routerName string) {
	fmt.Printf("# Table de routage RIP - %s\n", routerName)
	fmt.Printf("# Generated on %s\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Println()
	fmt.Println("# IP Destination | Masque de sous-réseau | Passerelle | Interface | Métrique")

	for _, interfaz := range interfaces {
		ip := interfaz.Interface.IP
		mask := cidrToSubnetMask(interfaz.Interface.Mascara)
		network := calculateNetworkAddress(ip, interfaz.Interface.Mascara)
		fmt.Printf("%s | %s | - | %s | 1\n", network, mask, interfaz.Interface.IP)
	}

	fmt.Println("# Fin de la table de routage.")
}

func main() {
	files := []string{
		"data/routeur-client.yaml",
		"data/routeur-r1.yaml",
		"data/routeur-r2.yaml",
		"data/routeur-r3.yaml",
		"data/routeur-r4.yaml",
		"data/routeur-r5.yaml",
		"data/routeur-r6.yaml",
		"data/routeur-serveur.yaml",
	}

	var buffer bytes.Buffer
	buffer.WriteByte(2)
	buffer.WriteByte(2)
	buffer.Write([]byte{0, 0})

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Error al leer el archivo %s: %v", file, err)
		}

		var interfaces []struct{ Interface Interfaz }
		err = yaml.Unmarshal(data, &interfaces)
		if err != nil {
			log.Fatalf("Error al deserializar el archivo %s: %v", file, err)
		}

		routerName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		printRIPFormat(interfaces, routerName)

		for _, interfaz := range interfaces {
			buffer.Write([]byte{0, 2})
			buffer.Write([]byte{0, 0})
			buffer.Write(ipToBytes(calculateNetworkAddress(interfaz.Interface.IP, interfaz.Interface.Mascara)))
			buffer.Write(maskToBytes(interfaz.Interface.Mascara))
			buffer.Write(ipToBytes("0.0.0.0"))
			buffer.Write([]byte{0, 0, 0, 1})
		}
	}

	err := ioutil.WriteFile("rip_message.bin", buffer.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error al escribir el archivo RIP: %v", err)
	}

	fmt.Println("Archivo RIP generado correctamente: rip_message.bin")
}
