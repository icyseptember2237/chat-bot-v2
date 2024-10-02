package config

import "chatbot/storage"

type Resource struct {
	Storage *storage.Storage `mapstructure:"storage"`
	Queue   *Queue           `mapstructure:"queue"`
}

type Queue struct {
}
