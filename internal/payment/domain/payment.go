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
	ID          string          `json:"ID"`
	CustomerID  string          `json:"CustomerID"`
	Status      string          `json:"Status"`
	PaymentLink string          `json:"PaymentLink"`
	Items       []*orderpb.Item `json:"Items"`
}
