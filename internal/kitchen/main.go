package main

import (
	"context"
	"github.com/xmhu2001/gorder-system/common/broker"
	grpcClient "github.com/xmhu2001/gorder-system/common/client"
	"github.com/xmhu2001/gorder-system/kitchen/adapters"
	"github.com/xmhu2001/gorder-system/kitchen/infrastructure/consumer"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/xmhu2001/gorder-system/common/config"
	"github.com/xmhu2001/gorder-system/common/logging"
	"github.com/xmhu2001/gorder-system/common/tracing"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("kitchen.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	orderClient, closeFn, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = closeFn()
	}()

	ch, closeConn := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeConn()
	}()

	orderGRPC := adapters.NewOrderGRPC(orderClient)
	go consumer.NewConsumer(orderGRPC).Listen(ch)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logrus.Info("receive signal, exiting...")
		os.Exit(0)
	}()
	logrus.Println("to exit, press ctrl+c")
	select {}
}
