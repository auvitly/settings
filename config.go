package settings

import (
	"io"

	"github.com/Auvitly/settings/internal"
	"github.com/Auvitly/settings/types"
	"github.com/spf13/viper"
)

// General

var configurator *internal.Configurator

// LoadOptions - функция загружает viper из файла
func LoadOptions(name string, paths ...string) (*viper.Viper, error) {
	configurator = internal.New(name, paths...)
	configurator.SetOption(types.LoggerHook, true)
	if err := configurator.LoadOptions(); err != nil {
		return nil, err
	}
	return configurator.Viper, nil
}

// LoadSettings - установка значения из viper в структуру, указатель на которую передается в качестве аргумента
func LoadSettings(settings interface{}, v *viper.Viper) error {
	return configurator.LoadSettings(settings)
}

// SetOption - настройка конфигуратора
func SetOption(options types.Options, value interface{}) error {
	return configurator.SetOption(options, value)
}

// New realisation

type IConfigurator interface {
	// ReadOptions - загрузка viper из внешнего reader
	ReadOptions(config io.Reader) error
	// LoadOptions - загрузка viper из файла, который был установлен при создании конфигуратора
	LoadOptions() error
	// LoadSettings - установка значения из viper в структуру, указатель на которую передается в качестве аргумента
	LoadSettings(config interface{}) error
	// SetOption - настройка конфигуратора
	SetOption(options types.Options, value interface{}) error
}

func New(name string, paths ...string) IConfigurator {
	return internal.New(name, paths...)
}
