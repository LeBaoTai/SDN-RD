package gnmib

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

func (c *Client) GetState(paths ...*gnmi.Path) (*gnmi.GetResponse, error) {
	ctx, cancel := c.ctxWithAuth()
	defer cancel()

	req := &gnmi.GetRequest{
		Path:     paths,
		Type:     gnmi.GetRequest_STATE,
		Encoding: gnmi.Encoding_JSON_IETF,
	}
	return c.gnmiC.Get(ctx, req)
}

func (c *Client) GetConfig(paths []*gnmi.Path) (*gnmi.GetResponse, error) {
	ctx, cancel := c.ctxWithAuth()
	defer cancel()

	req := &gnmi.GetRequest{
		Path:     paths,
		Type:     gnmi.GetRequest_CONFIG,
		Encoding: gnmi.Encoding_JSON_IETF,
	}
	return c.gnmiC.Get(ctx, req)
}

func (c *Client) GetConfigAsASCII(paths []*gnmi.Path) (*gnmi.GetResponse, error) {
	ctx, cancel := c.ctxWithAuth()
	defer cancel()

	req := &gnmi.GetRequest{
		Prefix: &gnmi.Path{
			Origin: "native",
		},
		Path:     paths,
		Type:     gnmi.GetRequest_CONFIG,
		Encoding: gnmi.Encoding_JSON_IETF,
	}
	return c.gnmiC.Get(ctx, req)
}

func (c *Client) GetOperational(paths ...*gnmi.Path) (*gnmi.GetResponse, error) {
	ctx, cancel := c.ctxWithAuth()
	defer cancel()

	req := &gnmi.GetRequest{
		Path:     paths,
		Type:     gnmi.GetRequest_OPERATIONAL,
		Encoding: gnmi.Encoding_JSON_IETF,
	}
	return c.gnmiC.Get(ctx, req)
}

func (c *Client) GetString(path string) (string, error) {
	p, err := BuildPath(path)
	if err != nil {
		return "", err
	}

	res, err := c.GetState(p)
	if err != nil {
		return "", err
	}

	for _, notif := range res.GetNotification() {
		for _, upd := range notif.GetUpdate() {
			if upd.GetVal() == nil {
				continue
			}

			jsonBytes := upd.GetVal().GetJsonIetfVal()
			if len(jsonBytes) == 0 {
				continue
			}

			return string(jsonBytes), nil
		}
	}
	return "", nil
}
