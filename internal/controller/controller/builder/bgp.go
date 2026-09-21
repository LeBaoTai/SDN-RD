package builder

type BGPReq struct {
	AS              string `json:"asn"`
	NeighborIP      string `json:"neighbor-ip"`
	PeerAS          string `json:"peer-as"`
	NetworkInstance string `json:"network-instance"`
}
