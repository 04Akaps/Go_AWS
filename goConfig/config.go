package goConfig

import (
	"log"

	"github.com/spf13/viper"
)

type GoConfig struct {
	AWS_REGION string `mapstructure:"aws_region"`
}

func LoadGoConfig(path string) GoConfig {
	var goConfig GoConfig

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("env Read Error : &w", err)
	}

	if err := viper.Unmarshal(&goConfig); err != nil {
		log.Fatal("env Marshal Error : &w", err)
	}

	return goConfig
}
