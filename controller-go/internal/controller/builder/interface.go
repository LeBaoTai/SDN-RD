package builder

import (
	"github.com/LeBaoTai/myco-controller/internal/oc"
	"github.com/openconfig/ygot/ygot"
)

type IfcReq struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
	Enable      bool   `json:"enable"`
}

func BuildInterface(ifcReq *IfcReq) {
	device := &oc.Device{}
	iface := device.GetOrCreateInterface(ifcReq.Name)
	iface.Description = ygot.String(ifcReq.Description)
	iface.Mtu = ygot.Uint16(1500)
	iface.AdminStatus = oc.Interface_AdminStatus_UP
}
