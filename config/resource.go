package config

import "chatbot/storage"

type Resource struct {
	Storage  *storage.Storage `mapstructure:"storage"`
	Chromium string           `mapstructure:"chromium"`
	Queue    *Queue           `mapstructure:"queue"`
}

type Queue struct {
}
