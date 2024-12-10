package domain

import (
	"context"

	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

// processor
type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}
