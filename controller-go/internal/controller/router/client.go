package router

import (
	"log"

	"github.com/golang/glog"
	"github.com/openconfig/gnmic/pkg/api/target"
	"github.com/openconfig/ygnmi/ygnmi"
)

type ClientCfg struct {
	Username string
	Password string
}

func CreateClient(t *target.Target) (*ygnmi.Client, error) {
	client, err := ygnmi.NewClient(
		t.Client,
		ygnmi.WithRequestLogLevel(glog.Level(log.Default().Flags())),
		ygnmi.WithTarget("router-core"),
	)
	if err != nil {
		return nil, err
	}
	log.Println(client.String())
	return client, nil
}
