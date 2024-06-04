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

// Estructura para representar la información de la interfaz
type Interfaz struct {
	Dispositivo string `yaml:"device"`
	IP          string `yaml:"ip"`
	Mascara     string `yaml:"mask"`
}

// Función para convertir una IP en un número binario
func ipToBytes(ip string) []byte {
	parsedIP := net.ParseIP(ip).To4()
	return parsedIP
}

// Función para convertir una máscara de subred en formato CIDR a bytes
func maskToBytes(mask string) []byte {
	var n int
	fmt.Sscanf(mask, "%d", &n)
	maskInt := ^uint32(0) << (32 - n)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, maskInt)
	return b
}

// Función para convertir una máscara de subred en formato CIDR a una notación de punto decimal
func cidrToSubnetMask(mask string) string {
	var n int
	fmt.Sscanf(mask, "%d", &n)
	maskInt := ^uint32(0) << (32 - n)
	return fmt.Sprintf("%d.%d.%d.%d", byte(maskInt>>24), byte(maskInt>>16), byte(maskInt>>8), byte(maskInt))
}

// Función para calcular la dirección de red a partir de la IP y la máscara de subred
func calculateNetworkAddress(ip, mask string) string {
	ipBytes := ipToBytes(ip)
	maskBytes := maskToBytes(mask)

	network := make([]byte, 4)
	for i := 0; i < 4; i++ {
		network[i] = ipBytes[i] & maskBytes[i]
	}

	return fmt.Sprintf("%d.%d.%d.%d", network[0], network[1], network[2], network[3])
}

// Función para imprimir la información de las interfaces en el formato RIP
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
	// Lista de archivos YAML para cada router
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

	// Añadir los campos fijos del encabezado RIP
	buffer.WriteByte(2)        // Comando: Respuesta
	buffer.WriteByte(2)        // Versión: RIP v2
	buffer.Write([]byte{0, 0}) // Dominio de enrutamiento: 0

	for _, file := range files {
		// Leer el archivo YAML
		data, err := ioutil.ReadFile(filepath.Join("..", file))
		if err != nil {
			log.Fatalf("Error al leer el archivo %s: %v", file, err)
		}

		// Deserializar el YAML en la estructura Interfaz
		var interfaces []struct{ Interface Interfaz }
		err = yaml.Unmarshal(data, &interfaces)
		if err != nil {
			log.Fatalf("Error al deserializar el archivo %s: %v", file, err)
		}

		// Obtener el nombre del router (suponiendo que el nombre del archivo es el nombre del router)
		routerName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		// Imprimir la información en el formato RIP
		printRIPFormat(interfaces, routerName)

		// Añadir las entradas de ruta
		for _, interfaz := range interfaces {
			buffer.Write([]byte{0, 2})                                                                          // Identifieur de famille d'@ : 2 (IPv4)
			buffer.Write([]byte{0, 0})                                                                          // Marqueur de route : 0
			buffer.Write(ipToBytes(calculateNetworkAddress(interfaz.Interface.IP, interfaz.Interface.Mascara))) // Adresse IP
			buffer.Write(maskToBytes(interfaz.Interface.Mascara))                                               // Masque de sous réseau
			buffer.Write(ipToBytes("0.0.0.0"))                                                                  // Passerelle : 0.0.0.0 (no gateway)
			buffer.Write([]byte{0, 0, 0, 1})                                                                    // Métrique : 1
		}
	}

	// Guardar el mensaje RIP en un archivo binario
	err := ioutil.WriteFile("../rip_message.bin", buffer.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error al escribir el archivo RIP: %v", err)
	}

	fmt.Println("Archivo RIP generado correctamente: rip_message.bin")
}
