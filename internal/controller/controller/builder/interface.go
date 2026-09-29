package builder

import (
	"context"
	"log"

	"github.com/LeBaoTai/SDN-RD/internal/oc"
	"github.com/LeBaoTai/SDN-RD/internal/oc/ocpath"
	"github.com/LeBaoTai/SDN-RD/internal/shared/model"
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

func CreateInterface(req *model.IfcReq) (*oc.Interface, error) {
	iface := &oc.Interface{}
	iface.Name = ygot.String(req.Name)
	iface.Description = ygot.String(req.Description)
	iface.Mtu = ygot.Uint16(req.Mtu)
	iface.Enabled = req.Enabled
	if req.Enabled != nil && *req.Enabled {
		iface.AdminStatus = oc.Interface_AdminStatus_UP
	} else {
		iface.AdminStatus = oc.Interface_AdminStatus_DOWN
	}
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

	// subindex
	var subIdx uint32 = 0
	if req.SubIndex != nil {
		subIdx = *req.SubIndex
	}

	// Subinterface
	subIface := iface.GetOrCreateSubinterface(*req.SubIndex)
	subIface.Enabled = req.Enabled
	subIface.Index = ygot.Uint32(subIdx)

	// IPV4
	ipv4 := subIface.GetOrCreateIpv4()
	ipv4.Enabled = ygot.Bool(true)

	addr := ipv4.GetOrCreateAddress(req.IP)
	addr.PrefixLength = ygot.Uint8(req.Mask)

	if err := iface.Validate(); err != nil {
		return nil, err
	}
	log.Printf("Valid Config for %v\n", req.Name)

	return iface, nil
}
