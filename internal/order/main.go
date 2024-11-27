package main

import (
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/config"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Printf("%v", viper.Get("order"))

}
