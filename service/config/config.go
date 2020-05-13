package config

import (
	"github.com/spf13/viper"

	"github.com/cherserver/micher/service/interfaces"
)

type config struct {
}

var _ interfaces.Config = &config{}

func New(configPath string) *config {
	viper.Get("key")
	return &config{}
}

func (c *config) Init() error {
	return nil
}
