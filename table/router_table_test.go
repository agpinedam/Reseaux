package table

import (
	"net"
	"testing"

	"reseaux/router"
)

func TestBuildRouteTable(t *testing.T) {
	// Definimos las interfaces del router para la prueba
	interfaces := []router.Interface{
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

	// Creamos una instancia del router con las interfaces definidas
	r := &router.Router{
		Interfaces: interfaces,
	}

	// Construimos la tabla de rutas usando la función BuildRouteTable
	routeTable := BuildRouteTable(r)

	// Verificamos que la tabla de rutas tenga el mismo número de interfaces
	if len(routeTable.Routes) != len(interfaces) {
		t.Fatalf("Expected %d routes, but got %d", len(interfaces), len(routeTable.Routes))
	}

	// Verificamos que cada interfaz en la tabla de rutas coincida con la configuración original
	for i, iface := range routeTable.Routes {
		if iface.Device != interfaces[i].Device || !iface.IP.Equal(interfaces[i].IP) || iface.Mask.String() != interfaces[i].Mask.String() {
			t.Errorf("Route %d does not match: expected %v, but got %v", i, interfaces[i], iface)
		}
	}
}
