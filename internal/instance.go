package internal

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"settings/types"
)

type Configurator struct {
	logger    *logrus.Logger
	fileName  string
	filePaths []string
	Viper     *viper.Viper
	config    map[string]interface{}
	options   map[types.Options]interface{}
	validator *validator.Validate
}

func New(name string, paths ...string) *Configurator {

	c := &Configurator{
		logger:    logrus.StandardLogger(),
		fileName:  name,
		filePaths: defaultPaths,
		Viper:     viper.New(),
		options:   make(map[types.Options]interface{}),
		validator: validator.New(),
	}

	// set filename
	c.Viper.SetConfigName(name)
	if len(name) == 0 {
		c.Viper.SetConfigName(defaultFileName)
	}

	// add path
	if len(paths) != 0 {
		c.filePaths = append(c.filePaths, paths...)
	}

	// if the filename is omitted, then use the default filename
	if len(name) == 0 {
		c.fileName = defaultFileName
	}

	// add base file paths
	for _, path := range c.filePaths {
		c.Viper.AddConfigPath(path)
	}

	// set default options
	c.setDefaultOptions()

	return c

}

func (c *Configurator) ReadConfiguration(config io.Reader) error {

	// loading settings into Viper
	err := c.Viper.ReadConfig(config)
	if err != nil {
		c.logger.WithError(err).Error("Unable to load file configuration from io.Reader")
		return err
	}

	return nil

}

func (c *Configurator) LoadConfiguration() error {

	// loading settings into Viper
	err := c.Viper.ReadInConfig()
	if err != nil {
		c.logger.WithError(err).Error("Unable to load file configuration from current paths")
		return err
	}
	c.config = c.Viper.AllSettings()

	return nil

}

func (c *Configurator) Unmarshal(config interface{}) error {

	if root, err := c.newRootHandler(config); err != nil {
		if c.getValidatorEnable() {
			if err = c.validator.Struct(root.reflectValue.Interface()); err != nil {
				return err
			}
		}
		return err
	} else {
		if err = c.handle(root); err != nil {
			return err
		}
		return nil
	}

}
