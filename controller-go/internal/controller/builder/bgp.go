package builder

import "github.com/LeBaoTai/myco-controller/internal/oc"

type BGPReq struct {
	AS              string `json:"asn"`
	NeighborIP      string `json:"neighbor-ip"`
	PeerAS          string `json:"peer-as"`
	NetworkInstance string `json:"network-instance"`
}

func (b *Builder) BuildBGP(req BGPReq) *oc.Device {
	//b.device.GetOrCreateComponent(Name string)
	return nil
}
