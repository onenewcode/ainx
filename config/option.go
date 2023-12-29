package config

import (
	"time"
)

// Option is the only struct that can be used to set Options.
type Option struct {
	F func(o *Options)
}

const (
	defaultKeepAliveTimeout = 1 * time.Minute
)

type Options struct {
	KeepAliveTimeout time.Duration
}

func (o *Options) Apply(opts []Option) {
	for _, op := range opts {
		op.F(o)
	}
}

func NewOptions(opts []Option) *Options {
	options := &Options{
		KeepAliveTimeout: defaultKeepAliveTimeout,
	}
	options.Apply(opts)
	return options
}
