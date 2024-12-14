package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/broker"
	_ "github.com/xmhu2001/gorder-system/common/config"
	"github.com/xmhu2001/gorder-system/common/discovery"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	"github.com/xmhu2001/gorder-system/common/logging"
	"github.com/xmhu2001/gorder-system/common/server"
	"github.com/xmhu2001/gorder-system/common/tracing"
	"github.com/xmhu2001/gorder-system/order/infrastructure/consumer"
	"github.com/xmhu2001/gorder-system/order/ports"
	"github.com/xmhu2001/gorder-system/order/service"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)
	// 在主函数里defer 关闭所有连接
	defer cleanup()

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
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
	go consumer.NewConsumer(application).Listen(ch)

	// 初始化消息队列 -> 移动到 application.go 了

	go server.RunGRPCServer(serviceName, func(s *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(s, svc)
	})
	// server 接受一个 app 作为胶水层，把handler等都粘合起来
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		// router.StaticFile 第一个参数相对路径指的是相对路由的路径
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
