package internal

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Options string

// General option types
const (
	TimeFormat     Options = "time_format"
	ProcessingMode Options = "processing_mode"
	LoggerHook     Options = "logger_hook"
	LoggerInstance Options = "logger_instance"
)

// ProcessingMode
const (
	OverwritingMode string = "overwriting"
	ComplementMode  string = "complement"
)

func (c *Configurator) SetOption(options Options, value interface{}) error {

	switch options {
	case TimeFormat:
		format, ok := value.(string)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[TimeFormat] = format
	case ProcessingMode:
		mode, ok := value.(string)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[ProcessingMode] = mode
	case LoggerHook:
		enable, ok := value.(bool)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[LoggerHook] = enable
	case LoggerInstance:
		logger, ok := value.(*logrus.Logger)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[LoggerHook] = logger
	default:
		return ErrInvalidOptionsType
	}

	return nil

}

func (c *Configurator) setDefaultOptions() {
	// Time format
	c.options[TimeFormat] = time.RFC3339
	// Processing mode
	c.options[ProcessingMode] = OverwritingMode
	// Logger hook
	c.options[LoggerHook] = false
	// Logger
	c.options[LoggerInstance] = logrus.StandardLogger()
}

func (c *Configurator) getTimeFormat() string {
	return c.options[TimeFormat].(string)
}

func (c *Configurator) getProcessingMode() string {
	return c.options[ProcessingMode].(string)
}

func (c *Configurator) getLoggerHook() bool {
	return c.options[LoggerHook].(bool)
}

func (c *Configurator) getLogger() *logrus.Logger {
	return c.options[LoggerInstance].(*logrus.Logger)
}
