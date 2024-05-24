package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Estructura para representar la información de la interfaz
type Interfaz struct {
	Dispositivo string `yaml:"device"`
	IP          string `yaml:"ip"`
	Mascara     string `yaml:"mask"`
}

func main() {
	// Leer el archivo YAML
	data, err := ioutil.ReadFile("../data/routeur-r5.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Deserializar el YAML en la estructura Interfaz
	var interfaces []struct{ Interface Interfaz }
	err = yaml.Unmarshal(data, &interfaces)
	if err != nil {
		log.Fatal(err)
	}

	// Imprimir la información de cada interfaz
	for _, interfaz := range interfaces {
		fmt.Printf("Dispositivo: %s\n", interfaz.Interface.Dispositivo)
		fmt.Printf("IP: %s\n", interfaz.Interface.IP)
		fmt.Printf("Mascara: %s\n\n", interfaz.Interface.Mascara)
	}
}
