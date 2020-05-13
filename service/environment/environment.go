package environment

import (
	"context"

	"github.com/cherserver/micher/service/interfaces"
)

type environment struct {
	shutdownCtx        context.Context
	shutdownCancelFunc context.CancelFunc

	config interfaces.Config
}

var _ interfaces.Environment = &environment{}

func New(config interfaces.Config) *environment {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &environment{
		shutdownCtx:        ctx,
		shutdownCancelFunc: cancelFunc,

		config: config,
	}
}

func (e *environment) Init() error {
	return nil
}

func (e *environment) ShutdownCtx() context.Context {
	return e.shutdownCtx
}
