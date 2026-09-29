package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/LeBaoTai/SDN-RD/internal/controller/controller/builder"
	"github.com/LeBaoTai/SDN-RD/internal/controller/controller/router"
	"github.com/LeBaoTai/SDN-RD/internal/controller/controller/utils"
	"github.com/LeBaoTai/SDN-RD/internal/shared/model"
)

func HandleInterfaceIntent(ctx context.Context, intent *model.IntentEnvelope, session *router.DeviceSession) error {
	var ifaceReq model.IfcReq
	err := json.Unmarshal(intent.Data, &ifaceReq)
	if err != nil {
		return err
	}
	iface, err := builder.CreateInterface(&ifaceReq)
	if err != nil {
		return err
	}

	authCtx := utils.WrapContextWithAuth(ctx, *session.Target.Config.Username, *session.Target.Config.Password)

	result, err := builder.UpdateInterface(iface, session.Client, authCtx)
	if err != nil {
		return err
	}
	log.Println(result)

	return nil
}
