package config

import (
	"settings/internal"
)

type Configurator interface {
}

func New(name string, path string) *internal.Configurator {
	return internal.New(name, path)
}
