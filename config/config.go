package config

import (
	"chatbot/function/config"
	"chatbot/logger"
)

var globalConfig = new(Config)

type Config struct {
	Log    logger.Options `mapstructure:"log"`
	Server config.Server  `mapstructure:"server"`

	Resources *Resource `mapstructure:"resources"`
}
