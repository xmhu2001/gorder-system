package order

import "github.com/xmhu2001/gorder-system/common/genproto/orderpb"

type Order struct {
	ID          string          `json:"id"`
	CustomerID  string          `json:"customer_id"`
	Status      string          `json:"status"`
	PaymentLink string          `json:"payment_link"`
	Items       []*orderpb.Item `json:"items"`
}
