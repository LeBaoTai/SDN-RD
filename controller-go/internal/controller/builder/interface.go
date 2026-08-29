package builder

import (
	"log"

	"github.com/LeBaoTai/myco-controller/internal/oc"
	"github.com/LeBaoTai/myco-controller/internal/oc/ocpath"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

type IfcReq struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
	Enabled     bool   `json:"enable"`
	SubIndex    uint32 `json:"sub-index"`
	Mask        uint8  `json:"mask"`
	Mtu         uint16 `json:"mtu"`
	Speed       int16  `json:"speed"`
	Duplex      string `json:"duplex"`
}

func CreateInterfaceUpdate(iface *oc.Interface) *gnmi.Update {
	ifacePath := ocpath.Root().Interface(*iface.Name)
	rawPath, _, err := ygnmi.ResolvePath(ifacePath)
	if err != nil {
		log.Printf("Cannot resolve path")
	}

	jsonBytes, err := ygot.EmitJSON(iface, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
	})
	if err != nil {
		log.Printf("Cannot create config: %v", err)
	}

	return &gnmi.Update{
		Path: rawPath,
		Val: &gnmi.TypedValue{
			Value: &gnmi.TypedValue_JsonIetfVal{
				JsonIetfVal: []byte(jsonBytes),
			},
		},
	}
}

func CreateInterface(req *IfcReq) *oc.Interface {
	iface := &oc.Interface{}

	iface.Name = ygot.String(req.Name)

	iface.Description = ygot.String(req.Description)
	iface.Mtu = ygot.Uint16(1500)
	iface.Enabled = ygot.Bool(req.Enabled)

	// ethernet
	eth := iface.GetOrCreateEthernet()
	switch req.Duplex {
	case "full":
		eth.DuplexMode = oc.Ethernet_DuplexMode_FULL
	case "auto":
		eth.DuplexMode = oc.Ethernet_DuplexMode_UNSET
	}

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
	subIface.Enabled = ygot.Bool(req.Enabled)

	// IPV4
	ipv4 := subIface.GetOrCreateIpv4()

	addr := ipv4.GetOrCreateAddress(req.IP)
	addr.PrefixLength = ygot.Uint8(req.Mask)

	return iface
}
