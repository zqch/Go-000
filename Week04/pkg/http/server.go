package http

import (
	"time"
)

type ServerOption struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int

	ShutdownTimeout time.Duration

	LogFunc func(v ...interface{})
}

func (t *ServerOption) Opt() ServeOptionFunc {
	return func(o *ServerOption) {
		if t.ReadTimeout != 0 {
			o.ReadTimeout = t.ReadTimeout
		}
		if t.ReadHeaderTimeout != 0 {
			o.ReadHeaderTimeout = t.ReadHeaderTimeout
		}
		if t.WriteTimeout != 0 {
			o.WriteTimeout = t.WriteTimeout
		}
		if t.IdleTimeout != 0 {
			o.IdleTimeout = t.IdleTimeout
		}
		if t.MaxHeaderBytes != 0 {
			o.MaxHeaderBytes = t.MaxHeaderBytes
		}
		if t.ShutdownTimeout != 0 {
			o.ShutdownTimeout = t.ShutdownTimeout
		}
		if t.LogFunc != nil {
			o.LogFunc = t.LogFunc
		}
	}
}

var defaultServerOption = ServerOption{
	ShutdownTimeout: 5 * time.Second,
	LogFunc:         println,
}

type ServeOptionFunc func(*ServerOption)
