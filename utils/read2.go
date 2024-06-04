//package utils
/*
import (
	"fmt"
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v2"
)

// Estructura para representar la información de la interfaz
type Interfaz struct {
	Dispositivo string `yaml:"device"`
	IP          string `yaml:"ip"`
	Mascara     string `yaml:"mask"`
}

// Leer el archivo YAML y convertir a octetos RIP
func ReadYAMLToRIP(filePath string) ([]byte, error) {
	// Leer el archivo YAML
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Deserializar el YAML en la estructura Interfaz
	var interfaces []struct{ Interface Interfaz }
	err = yaml.Unmarshal(data, &interfaces)
	if err != nil {
		return nil, err
	}

	// Convertir a formato RIP en octetos
	var ripData []byte
	for _, interfaz := range interfaces {
		ip := net.ParseIP(interfaz.Interface.IP).To4()
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", interfaz.Interface.IP)
		}
		mask := net.ParseIP(interfaz.Interface.Mascara).To4()
		if mask == nil {
			return nil, fmt.Errorf("invalid subnet mask: %s", interfaz.Interface.Mascara)
		}
		device := []byte(interfaz.Interface.Dispositivo)
		device = append(device, make([]byte, 16-len(device))...) // Rellenar a 16 bytes
		ripData = append(ripData, device...)
		ripData = append(ripData, ip...)
		ripData = append(ripData, mask...)
	}
	return ripData, nil
}
*/