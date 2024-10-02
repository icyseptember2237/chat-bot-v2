package config

import "chatbot/logger"

type Server struct {
	Port string          `mapstructure:"port"`
	Log  *logger.Options `mapstructure:"log"`
}
