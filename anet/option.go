package anet

import (
	"ainx/config"
	"time"
)

func WithKeepAliveTimeout(t time.Duration) config.Option {
	return config.Option{F: func(o *config.Options) {
		o.KeepAliveTimeout = t
	}}
}
