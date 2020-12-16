package http

import (
	"time"
)

type ServerOption struct {
	ShutdownTimeout time.Duration
	LogFunc func(v ...interface{})
}

var defaultServerOption = ServerOption{
	ShutdownTimeout: 5 * time.Second,
	LogFunc:         println,
}

type ServeOptionFunc func(*ServerOption)

func (t *ServerOption) Opt() ServeOptionFunc {
	return func(o *ServerOption) {
		if t.ShutdownTimeout != 0 {
			o.ShutdownTimeout = t.ShutdownTimeout
		}
		if t.LogFunc != nil {
			o.LogFunc = t.LogFunc
		}
	}
}

func New() {

}