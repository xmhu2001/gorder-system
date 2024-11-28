package main

import (
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/genproto/stockpb"
	"github.com/xmhu2001/gorder-system/common/server"
	"github.com/xmhu2001/gorder-system/stock/ports"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("grpc")
	serverType := viper.GetString("stock.server-to-run")
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(s *grpc.Server) {
			stockpb.RegisterStockServiceServer(s, ports.NewGRPCServer())
		})
	case "http":
		// TODO
	default:
		panic("unknown server type")
	}

}
