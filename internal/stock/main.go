package main

import (
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/server"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("grpc")
	server.RunGRPCServer(serviceName, func(s *grpc.Server) {

	})
}
