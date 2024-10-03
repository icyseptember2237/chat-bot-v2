package config

import (
	"chatbot/function/config"
	config2 "chatbot/job/config"
	"chatbot/logger"
)

var globalConfig = new(Config)

type Config struct {
	Log    logger.Options      `mapstructure:"log"`
	Server config.Server       `mapstructure:"server"`
	Jobs   []config2.JobConfig `mapstructure:"jobs"`

	Resources *Resource `mapstructure:"resources"`
}
