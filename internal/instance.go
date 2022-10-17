package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configurator struct {
	logger    *logrus.Logger
	fileName  string
	filePaths []string
	viper     *viper.Viper
}

func New(name string, path string) *Configurator {

	c := &Configurator{
		logger:    logrus.StandardLogger(),
		fileName:  name,
		filePaths: defaultPaths,
		viper:     viper.New(),
	}

	// set filename
	c.viper.SetConfigName(name)
	if len(name) == 0 {
		c.viper.SetConfigName(defaultFileName)
	}

	// add path
	if len(path) != 0 {
		c.filePaths = append(c.filePaths, path)
	}

	// if the filename is omitted, then use the default filename
	if len(name) == 0 {
		c.fileName = defaultFileName
	}

	// add base file paths
	for _, path = range c.filePaths {
		c.viper.AddConfigPath(path)
	}

	return c

}

func (c *Configurator) AddConfigPath(paths ...string) {
	c.filePaths = append(c.filePaths, paths...)
	for _, path := range paths {
		c.viper.AddConfigPath(path)
	}
}

func (c *Configurator) LoadConfiguration() error {

	// loading settings into viper
	err := c.viper.ReadInConfig()
	if err != nil {
		c.logger.WithError(err).Error("Unable to load file configuration from current paths")
		return err
	}

	return nil

}

func (c *Configurator) UnmarshalInto(config interface{}) error {
	return c.handle(config)

}
