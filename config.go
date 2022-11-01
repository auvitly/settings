package config

import (
	"github.com/spf13/viper"
	"io"
	"settings/internal"
)

type IConfigurator interface {
	ReadConfiguration(config io.Reader) error
	LoadConfiguration() error
	Unmarshal(config interface{}) error
	SetOption(options internal.IOptions, value interface{}) error
	GetViper() (*viper.Viper, error)
	AddFilePaths(paths ...string)
}

func New(name string, path string) *internal.Configurator {
	return internal.New(name, path)
}

// Old realisation

var configurator IConfigurator

func LoadOptions(name string, paths ...string) (*viper.Viper, error) {
	configurator = internal.New(name, "")
	configurator.SetOption(LoggerHook, true)
	if len(paths) != 0 {
		configurator.AddFilePaths(paths...)
	}
	if err := configurator.LoadConfiguration(); err != nil {
		return nil, err
	}
	return configurator.GetViper()
}

func LoadSettings(settings interface{}, v *viper.Viper) error {
	return configurator.Unmarshal(settings)
}
