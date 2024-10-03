package config

import (
	"chatbot/function/config"
	jobconfig "chatbot/job/config"
	"chatbot/logger"
	workerconfig "chatbot/worker/config"
)

var globalConfig = new(Config)

type Config struct {
	Log     logger.Options              `mapstructure:"log"`
	Server  config.Server               `mapstructure:"server"`
	Jobs    []jobconfig.JobConfig       `mapstructure:"jobs"`
	Workers []workerconfig.WorkerConfig `mapstructure:"workers"`

	Resources *Resource `mapstructure:"resources"`
}
