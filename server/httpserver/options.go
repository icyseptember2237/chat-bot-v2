package httpserver

import (
	"gotemplate/server/httpserver/middles"
)

type Options struct {
	Name    string
	Address string
	Middles []middles.Middle
}

func newOptions(opts ...Option) Options {
	options := Options{
		Name:    "httpserver",
		Address: ":8080",
		Middles: nil,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type Option func(*Options)

func WithName(name string) Option {
	return func(opt *Options) {
		opt.Name = name
	}
}

func WithAddress(address string) Option {
	return func(opt *Options) {
		opt.Address = address
	}
}

func WithMiddles(ms ...middles.Middle) Option {
	return func(opt *Options) {
		opt.Middles = ms
	}
}
