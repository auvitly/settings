package config

import (
	"github.com/spf13/viper"
	"io"
	"settings/internal"
	"settings/types"
)

// New realisation

type IConfigurator interface {
	// ReadConfiguration - загрузка viper из внешнего reader
	ReadConfiguration(config io.Reader) error
	// LoadConfiguration - загрузка загрузить viper из файла, который был установлен при создании конфигуратора
	LoadConfiguration() error
	// Unmarshal - установка значения из viper в структуру, указатель на которую передается в качестве аргумента
	Unmarshal(config interface{}) error
	// SetOption - настройка конфигуратора
	SetOption(options types.Options, value interface{}) error
	// GetViper - возвращает *viper.Viper конфигуратора
	GetViper() (*viper.Viper, error)
}

func New(name string, paths ...string) *internal.Configurator {
	return internal.New(name, paths...)
}

// Old realisation

var configurator IConfigurator

// LoadOptions - функция загружает viper из файла
func LoadOptions(name string, paths ...string) (*viper.Viper, error) {
	configurator = internal.New(name, paths...)
	configurator.SetOption(types.LoggerHook, true)
	if err := configurator.LoadConfiguration(); err != nil {
		return nil, err
	}
	return configurator.GetViper()
}

// LoadSettings - установка значения из viper в структуру, указатель на которую передается в качестве аргумента
func LoadSettings(settings interface{}, v *viper.Viper) error {
	return configurator.Unmarshal(settings)
}
