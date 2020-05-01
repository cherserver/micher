package config

import "github.com/cherserver/micher/service/interfaces"

type config struct {
}

var _ interfaces.Config = config{}

func New() *config {
	return &config{}
}
