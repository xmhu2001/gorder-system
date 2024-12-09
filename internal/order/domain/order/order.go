package order

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v81"
	"github.com/xmhu2001/gorder-system/order/entity"
)

// Aggregate
type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

// 业务逻辑写在domain
func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item) (*Order, error) {
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

func (o *Order) IsPaid() error {
	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
		return nil
	}
	return fmt.Errorf("order status not paid, order status=%s, order id=%s", o.Status, o.ID)
}
