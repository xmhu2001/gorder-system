package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	grpcClient "github.com/xmhu2001/gorder-system/common/client"
	"github.com/xmhu2001/gorder-system/common/metrics"
	"github.com/xmhu2001/gorder-system/payment/adapters"
	"github.com/xmhu2001/gorder-system/payment/app"
	"github.com/xmhu2001/gorder-system/payment/app/command"
	"github.com/xmhu2001/gorder-system/payment/domain"
	"github.com/xmhu2001/gorder-system/payment/infrastructure/processor"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	// memoryProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
