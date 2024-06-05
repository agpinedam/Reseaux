package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"reseaux/utils" // Ajusta esta ruta según tu estructura de proyecto
)

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

func printRIPFormat(interfaces []struct{ Interface utils.Interfaz }, routerName string) {
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

func generateRIPMessage(files []string) ([]byte, error) {
	dataMap, err := utils.ReadYAMLFiles(files)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.WriteByte(2)        // Comando: Respuesta
	buffer.WriteByte(2)        // Versión: RIP v2
	buffer.Write([]byte{0, 0}) // Dominio de enrutamiento: 0

	for routerName, interfaces := range dataMap {
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
	return buffer.Bytes(), nil
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

	ripMessage, err := generateRIPMessage(files)
	if err != nil {
		log.Fatalf("Error al generar el mensaje RIP: %v", err)
	}

	err = sendRIPMessage(ripMessage, "localhost:8080")
	if err != nil {
		log.Fatalf("Error al enviar el mensaje RIP: %v", err)
	}

	fmt.Println("Mensaje RIP enviado correctamente")
}

func sendRIPMessage(message []byte, address string) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return fmt.Errorf("error conectando al servidor: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		return fmt.Errorf("error enviando el mensaje: %v", err)
	}
	return nil
}
