package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	EnabledChats []string `mapstructure:"enabled_chats"`
	LogLevel string `mapstructure:"log_level"`
	Twitch TwitchConfig `mapstructure:"twitch"`
}

type TwitchConfig struct {
	Channel string `mapstructure:"channel"`
}

var conf = Config{}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/") // Container Environment
	viper.AddConfigPath(".") // Local Environment

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}

	logLevel, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		panic(err)
	}

	logrus.SetLevel(logLevel)
	if logLevel != logrus.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}
}

func GetConfig() *Config {
	return &conf
}