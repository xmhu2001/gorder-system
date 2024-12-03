package order

import (
	"errors"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

type Order struct {
	ID          string          `json:"id"`
	CustomerID  string          `json:"CustomerID"`
	Status      string          `json:"Status"`
	PaymentLink string          `json:"PaymentLink"`
	Items       []*orderpb.Item `json:"Items"`
}

// 业务逻辑写在domain
func NewOrder(id, customerID, status, paymentLink string, items []*orderpb.Item) (*Order, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if customerID == "" {
		return nil, errors.New("customerID is required")
	}
	if status == "" {
		return nil, errors.New("status is required")
	}
	if items == nil {
		return nil, errors.New("items is required")
	}
	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      status,
		PaymentLink: paymentLink,
		Items:       items,
	}, nil
}

func (o *Order) ToProto() *orderpb.Order {
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       o.Items,
		PaymentLink: o.PaymentLink,
	}
}
