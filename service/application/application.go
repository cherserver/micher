package application

import "github.com/cherserver/micher/service/interfaces"

type application struct {
	env interfaces.Environment
}

func New(environment interfaces.Environment) *application {
	return &application{
		env: environment,
	}
}

func (a *application) Init() error {
	return nil
}
