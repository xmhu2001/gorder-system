package processor

import (
	"context"
	"encoding/json"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

type StripeProcessor struct {
	apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	if apiKey == "" {
		panic("empty apiKey")
	}
	stripe.Key = apiKey
	return &StripeProcessor{apiKey: apiKey}
}

var (
	successURL = "http://localhost:8080/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	// https://docs.stripe.com/payments/checkout/how-checkout-works?lang=go
	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			// item.PriceID
			Price:    stripe.String("price_1QRSEwEnyEnePsWo9ITjIIZ8"),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	marshalledItems, _ := json.Marshal(order.Items)
	metadata := map[string]string{
		"orderID":    order.ID,
		"customerID": order.CustomerID,
		"status":     order.Status,
		"items":      string(marshalledItems),
	}
	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
