package config

import (
	"io"
	"settings/internal"
)

type Configurator interface {
	ReadConfiguration(config io.Reader) error
	LoadConfiguration() error
	Unmarshal(config interface{}) error
	SetOption(options internal.Options, value interface{})
}

func New(name string, path string) *internal.Configurator {
	return internal.New(name, path)
}
