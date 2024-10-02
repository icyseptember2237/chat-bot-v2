package config

var globalConfig = new(Config)

type Config struct {
	Server    *Server   `mapstructure:"server,omitempty"`
	Resources *Resource `mapstructure:"resources"`
}
