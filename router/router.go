package router

import (
	"fmt"
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v2"
)

type Router struct {
	Interfaces []Interface `yaml:"interfaces"`
}

type Interface struct {
	Device string `yaml:"device"`
	IP     net.IP `yaml:"ip"`
	Mask   net.IPMask
	Metric int // Agrega el campo para la métrica
}

type rawInterface struct {
	Device string `yaml:"device"`
	IP     string `yaml:"ip"`
	Mask   int    `yaml:"mask"`
}

func NewRouterFromFile(filePath string) (*Router, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var rawInterfaces struct {
		Interfaces []rawInterface `yaml:"interfaces"`
	}
	if err := yaml.Unmarshal(data, &rawInterfaces); err != nil {
		return nil, err
	}

	var interfaces []Interface
	for _, rawIface := range rawInterfaces.Interfaces {
		ip := net.ParseIP(rawIface.IP)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", rawIface.IP)
		}

		mask := net.CIDRMask(rawIface.Mask, 32)
		iface := Interface{
			Device: rawIface.Device,
			IP:     ip,
			Mask:   mask,
			Metric: 1, // Inicializa la métrica en 1
		}
		interfaces = append(interfaces, iface)
	}

	return &Router{Interfaces: interfaces}, nil
}
