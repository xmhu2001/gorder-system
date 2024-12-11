package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

var once sync.Once

// 单例模式保护 NewViperConfig()
func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return
}

func newViperConfig() error {
	realPath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(realPath)
	// viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	return viper.ReadInConfig()
}

func getRelativePathFromCaller() (realPath string, err error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return
	}
	_, here, _, _ := runtime.Caller(0)
	realPath, err = filepath.Rel(callerPwd, filepath.Dir(here))
	fmt.Printf("realpath: %s callerpath: %s curpath: %s\n", realPath, callerPwd, here)
	return
}
