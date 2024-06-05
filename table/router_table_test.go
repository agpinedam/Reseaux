package table

import (
	"net"
	"testing"

	"reseaux/router"
)

func TestBuildRouteTable(t *testing.T) {
	// Ruta al archivo YAML de prueba
	routerConfigPath := "../data/routeur-r1.yaml"

	// Cargamos el router desde el archivo YAML
	r, err := router.NewRouterFromFile(routerConfigPath)
	if err != nil {
		t.Fatalf("Failed to load router configuration: %v", err)
	}

	// Construimos la tabla de rutas usando la función BuildRouteTable
	routeTable := BuildRouteTable(r)

	// Definimos las interfaces esperadas
	expectedInterfaces := []router.Interface{
		{
			Device: "eth0",
			IP:     net.ParseIP("192.168.1.254"),
			Mask:   net.CIDRMask(24, 32),
		},
		{
			Device: "eth1",
			IP:     net.ParseIP("10.1.1.1"),
			Mask:   net.CIDRMask(30, 32),
		},
	}

	// Verificamos que la tabla de rutas tenga el mismo número de interfaces
	if len(routeTable.Routes) != len(expectedInterfaces) {
		t.Fatalf("Expected %d routes, but got %d", len(expectedInterfaces), len(routeTable.Routes))
	}

	// Verificamos que cada interfaz en la tabla de rutas coincida con la configuración esperada
	for i, iface := range routeTable.Routes {
		expectedIface := expectedInterfaces[i]
		if iface.Device != expectedIface.Device || !iface.IP.Equal(expectedIface.IP) || iface.Mask.String() != expectedIface.Mask.String() {
			t.Errorf("Route %d does not match: expected %v, but got %v", i, expectedIface, iface)
		}
	}
}
