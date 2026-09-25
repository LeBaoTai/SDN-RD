package builder

import (
	"context"
	"log"

	"github.com/LeBaoTai/SDN-RD/internal/controller/nats"
	"github.com/LeBaoTai/SDN-RD/internal/oc"
	"github.com/LeBaoTai/SDN-RD/internal/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

func UpdateInterface(iface *oc.Interface, client *ygnmi.Client, ctx context.Context) (*ygnmi.Result, error) {
	ifaceConfigQuery := ocpath.Root().Interface(*iface.Name).Config()

	result, err := ygnmi.Update(ctx, client, ifaceConfigQuery, iface)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CreateInterface(req *nats.IfcReq) (*oc.Interface, error) {
	iface := &oc.Interface{}
	iface.Name = new(req.Name)
	iface.Description = new(req.Description)
	iface.Mtu = ygot.Uint16(1500)
	iface.Enabled = new(req.Enabled)
	iface.Type = oc.IETFInterfaces_InterfaceType_ethernetCsmacd

	// ethernet
	eth := iface.GetOrCreateEthernet()
	switch req.Speed {
	case 100:
		eth.PortSpeed = oc.OpenconfigIfEthernet_ETHERNET_SPEED_SPEED_100MB
	case 1000:
		eth.PortSpeed = oc.OpenconfigIfEthernet_ETHERNET_SPEED_SPEED_1GB
	case 10000:
		eth.PortSpeed = oc.OpenconfigIfEthernet_ETHERNET_SPEED_SPEED_10GB
	}

	// Subinterface
	subIface := iface.GetOrCreateSubinterface(req.SubIndex)
	subIface.Enabled = new(req.Enabled)
	subIface.Index = new(uint32(0))

	// IPV4
	ipv4 := subIface.GetOrCreateIpv4()
	ipv4.Enabled = new(true)

	addr := ipv4.GetOrCreateAddress(req.IP)
	addr.PrefixLength = new(req.Mask)

	if err := iface.Validate(); err != nil {
		return nil, err
	}
	log.Printf("Valid Config for %v\n", req.Name)

	return iface, nil
}
