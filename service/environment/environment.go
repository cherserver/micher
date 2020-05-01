package environment

import "github.com/cherserver/micher/service/interfaces"

type environment struct {
}

var _ interfaces.Environment = environment{}

func New() *environment {
	return &environment{}
}

func (e *environment) Init() error {
	return nil
}
