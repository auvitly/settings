package internal

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"settings/types"
)

type Configurator struct {
	logger    *logrus.Logger
	fileName  string
	filePaths []string
	viper     *viper.Viper
	config    map[string]interface{}
	options   map[types.Options]interface{}
	validator *validator.Validate
}

func New(name string, paths ...string) *Configurator {

	c := &Configurator{
		logger:    logrus.StandardLogger(),
		fileName:  name,
		filePaths: defaultPaths,
		viper:     viper.New(),
		options:   make(map[types.Options]interface{}),
		validator: validator.New(),
	}

	// set filename
	c.viper.SetConfigName(name)
	if len(name) == 0 {
		c.viper.SetConfigName(defaultFileName)
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
		c.viper.AddConfigPath(path)
	}

	// set default options
	c.setDefaultOptions()

	return c

}

func (c *Configurator) ReadConfiguration(config io.Reader) error {

	// loading settings into viper
	err := c.viper.ReadConfig(config)
	if err != nil {
		c.logger.WithError(err).Error("Unable to load file configuration from io.Reader")
		return err
	}

	return nil

}

func (c *Configurator) LoadConfiguration() error {

	// loading settings into viper
	err := c.viper.ReadInConfig()
	if err != nil {
		c.logger.WithError(err).Error("Unable to load file configuration from current paths")
		return err
	}
	c.config = c.viper.AllSettings()

	return nil

}

func (c *Configurator) Unmarshal(config interface{}) error {

	if root, err := c.newRootHandler(config); err != nil {
		if err = c.validator.Struct(root.reflectValue.Interface()); err != nil {
			return err
		}
		return err
	} else {
		if err = c.handle(root); err != nil {
			return err
		}
		return nil
	}

}

func (c *Configurator) GetViper() (*viper.Viper, error) {
	if c.viper != nil {
		return c.viper, nil
	} else {
		return nil, errors.New("viper not found")
	}
}
