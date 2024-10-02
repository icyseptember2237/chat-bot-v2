package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gotemplate/logger"
	"reflect"
)

func NewLoggerHook() Hook {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() != reflect.Map {
			return data, nil
		}
		if to == reflect.TypeOf(logger.Options{}) {
			var options logger.Options
			err := mapstructure.Decode(data, &options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode data to logger.Standard")
			}
			err = logger.ResetStandardWithOptions(options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to reset standard")
			}
			return options, nil
		}
		return data, nil
	}
}
