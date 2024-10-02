package config

import "gotemplate/logger"

type Server struct {
	Port string          `mapstructure:"port"`
	Log  *logger.Options `mapstructure:"log"`
}
