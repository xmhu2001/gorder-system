package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/config"
	"github.com/xmhu2001/gorder-system/common/server"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// todo
	serviceName := viper.GetString("order.service-name")
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {

	})
}
