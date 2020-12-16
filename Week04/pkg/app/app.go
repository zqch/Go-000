package app

import "context"

type App interface {
	Run(context.Context) error
}
