package interfaces

import "context"

type Environment interface {
	Init() error

	ShutdownCtx() context.Context
}
