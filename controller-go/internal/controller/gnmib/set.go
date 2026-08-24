package gnmib

import (
	"context"

	"github.com/openconfig/gnmi/proto/gnmi"
)

type Transaction struct {
	ctx context.Context
	sr  *gnmi.SetRequest
}

func (t *TargetClient) CreateTransaction(ctx context.Context, deletes []*gnmi.Path, updates []*gnmi.Update, replaces []*gnmi.Update) *Transaction {
	sr := &gnmi.SetRequest{
		Delete:  deletes,
		Update:  updates,
		Replace: replaces,
	}
	return &Transaction{
		ctx: ctx,
		sr:  sr,
	}
}

func (t *TargetClient) SetTransaction(tr *Transaction) (*gnmi.SetResponse, error) {
	return t.target.Set(tr.ctx, tr.sr)
}
