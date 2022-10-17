package config

import (
	"settings/internal"
)

type Configurator interface {
	LoadConfiguration() error
	Unmarshal(config interface{}) error
}

func New(name string, path string) *internal.Configurator {
	return internal.New(name, path)
}
