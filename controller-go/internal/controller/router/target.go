package router

import (
	"context"
	"time"

	"github.com/openconfig/gnmic/pkg/api"
	"github.com/openconfig/gnmic/pkg/api/target"
)

type TargetCfg struct {
	Address    string
	Port       string
	SkipVerify bool
	Timeout    time.Duration
}

func CreateNewTarget(cfg TargetCfg, ctx context.Context) (*target.Target, error) {
	target, err := api.NewTarget(
		api.Address(cfg.Address+":"+cfg.Port),
		api.SkipVerify(cfg.SkipVerify),
		api.Timeout(cfg.Timeout),
	)
	if err != nil {
		return nil, err
	}
	return target, nil
}
