package golinks

import (
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Port string

	Storage storageConfig
}

type storageConfig struct {
	Type string
}

var C config

func init() {
	viper.SetEnvPrefix("golinks")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	viper.SetDefault("Port", 8080)
	viper.BindEnv("Port")
	viper.SetDefault("Storage.Type", "memory")
	viper.BindEnv("Storage.Type")

	err := viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
}
