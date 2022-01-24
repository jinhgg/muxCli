package config

import "github.com/spf13/viper"

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	_ = viper.ReadInConfig()
}
