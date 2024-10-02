package config

import "gotemplate/storage"

type Resource struct {
	Storage *storage.Storage `mapstructure:"storage"`
	Queue   *Queue           `mapstructure:"queue"`
}

type Queue struct {
}
