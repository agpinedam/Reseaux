package utils

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Interfaz struct {
	Dispositivo string `yaml:"device"`
	IP          string `yaml:"ip"`
	Mascara     string `yaml:"mask"`
}

// Function to read and parse YAML files
func ReadYAMLFiles(files []string) (map[string][]struct{ Interface Interfaz }, error) {
	dataMap := make(map[string][]struct{ Interface Interfaz })

	for _, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(file))
		if err != nil {
			log.Fatalf("Error al leer el archivo %s: %v", file, err)
			return nil, err
		}

		var interfaces []struct{ Interface Interfaz }
		err = yaml.Unmarshal(data, &interfaces)
		if err != nil {
			log.Fatalf("Error al deserializar el archivo %s: %v", file, err)
			return nil, err
		}

		routerName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		dataMap[routerName] = interfaces
	}

	return dataMap, nil
}
