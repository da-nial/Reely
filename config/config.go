package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatalln("error on parsing configuration file", err)
	}
}

func GetConfig() *viper.Viper {
	return config
}
