package dynamodb

import (
	"github.com/spf13/viper"
)

type config struct {
	Storage storageConfig
}

type storageConfig struct {
	Region    string
	Endpoint  string
	TableName string
}

var C config

func loadConfig() {
	viper.SetDefault("Storage.Region", "us-west-2")
	viper.BindEnv("Storage.Region")

	viper.BindEnv("Storage.Endpoint")

	viper.SetDefault("Storage.TableName", "Links")
	viper.BindEnv("Storage.TableName")

	err := viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
}
