package internal

import (
	"io"

	"github.com/Auvitly/settings/types"
	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/sirupsen/logrus"
)

func (c *Configurator) configureLogger(config types.Logger) error {

	logger := c.getLogger()

	if !config.StdOut {
		logger.Out = io.Discard
	}

	logger.Formatter = &logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    !config.Colour,
	}

	logger.SetLevel(config.LogLevel)

	// активируем хук в сислог
	hook, err := syslogHook(config)
	if err != nil {
		return err
	} else if hook != nil {
		logger.AddHook(hook)
	}

	// активируем хук в грейлог
	if len(config.Graylog) != 0 {
		graylogHook := graylog.NewAsyncGraylogHook(config.Graylog, nil)
		graylogHook.Level = config.GraylogLevel
		logger.AddHook(graylogHook)
	}

	return nil
}
