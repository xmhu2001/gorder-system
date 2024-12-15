package integration

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v81/product"

	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v81"
	_ "github.com/xmhu2001/gorder-system/common/config"
)

type StripeAPI struct {
	apiKey string
}

func NewStripeAPI() *StripeAPI {
	key := viper.GetString("stripe-key")
	if key == "" {
		logrus.Fatal("empty key")
	}
	return &StripeAPI{apiKey: key}
}

func (s *StripeAPI) GetPriceByProductID(ctx context.Context, productID string) (string, error) {
	// TODO: log...
	stripe.Key = s.apiKey
	result, err := product.Get(productID, &stripe.ProductParams{})
	if err != nil {
		return "", err
	}
	return result.DefaultPrice.ID, nil
}
