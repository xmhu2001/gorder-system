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

func getRelativePathFromCaller() (relPath string, err error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return
	}
	_, here, _, _ := runtime.Caller(0)
	relPath, err = filepath.Rel(callerPwd, filepath.Dir(here))
	fmt.Printf("relpath: %s callerpath: %s curpath: %s\n", relPath, callerPwd, here)
	return
}
