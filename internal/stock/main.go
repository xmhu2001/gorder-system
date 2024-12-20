package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/xmhu2001/gorder-system/common/config"
	"github.com/xmhu2001/gorder-system/common/discovery"
	"github.com/xmhu2001/gorder-system/common/genproto/stockpb"
	"github.com/xmhu2001/gorder-system/common/logging"
	"github.com/xmhu2001/gorder-system/common/server"
	"github.com/xmhu2001/gorder-system/common/tracing"
	"github.com/xmhu2001/gorder-system/stock/ports"
	"github.com/xmhu2001/gorder-system/stock/service"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application := service.NewApplication(ctx)
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(s *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(s, svc)
		})
	case "http":
		// TODO
	default:
		panic("unknown server type")
	}

}
