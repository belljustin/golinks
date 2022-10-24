package dynamodb

import (
	"github.com/spf13/viper"
)

type config struct {
	Storage storageConfig
}

type storageConfig struct {
	Region   string
	Endpoint string
}

var C config

func loadConfig() {
	viper.SetDefault("Storage.Region", "us-west-2")
	viper.BindEnv("Storage.Region")
	viper.SetDefault("Storage.Endpoint", "http://localhost:8000")
	viper.BindEnv("Storage.Endpoint")

	err := viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
}
