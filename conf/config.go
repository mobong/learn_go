package conf

import (
	"github.com/spf13/viper"
)

func ConfInit() {
	viper.SetConfigType("toml")
	viper.SetConfigFile("conf/config.toml")
	viper.ReadInConfig()
}
