// En el archivo router/router.go
package router

import (
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v2"
)

type Router struct {
	Interfaces []Interface `yaml:"interface"`
}

type Interface struct {
	Device string     `yaml:"device"`
	IP     net.IP     `yaml:"ip"`
	Mask   net.IPMask `yaml:"mask"`
}

func NewRouterFromFile(filePath string) (*Router, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var router Router
	if err := yaml.Unmarshal(data, &router); err != nil {
		return nil, err
	}

	return &router, nil
}
