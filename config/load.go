package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Hook = mapstructure.DecodeHookFunc

var defaultHooks = []Hook{
	NewLoggerHook(),
}

func Load(path string) (*Config, error) {
	err := LoadWithHooks(path, globalConfig, defaultHooks...)
	if err != nil {
	}
	return globalConfig, err
}

func LoadWithHooks(path string, conf interface{}, hooks ...Hook) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	err := viper.Unmarshal(conf, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(hooks...)))
	if err != nil {
		return errors.Wrap(err, "failed to scan config")
	}
	return nil
}

func Get() *Config {
	return globalConfig
}
