package main

import (
	"context"
	"github.com/xmhu2001/gorder-system/common/tracing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/broker"
	_ "github.com/xmhu2001/gorder-system/common/config"
	"github.com/xmhu2001/gorder-system/common/logging"
	"github.com/xmhu2001/gorder-system/common/server"
	"github.com/xmhu2001/gorder-system/payment/infrastructure/consumer"
	"github.com/xmhu2001/gorder-system/payment/service"
)

func init() {
	logging.Init()
}

// 1. consumer listen EventOrderCreated 事件 （consumer.go）
// 2. 消费这个消息后, 去 handleMessage(), 在这里创建支付链接

func main() {
	serviceName := viper.GetString("payment.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serverType := viper.GetString("payment.server-to-run")

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	// 初始化消息队列
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

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported type grpc") // 还会一层一层执行调用栈内地一些defer函数，Fatal则是直接退出
	default:
		logrus.Panic("unsupported type")
	}
}
