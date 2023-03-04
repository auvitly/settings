package internal

import (
	"time"

	"github.com/Auvitly/settings/types"
	"github.com/sirupsen/logrus"
)

func (c *Configurator) SetOption(options types.Options, value interface{}) error {

	switch options {
	case types.TimeFormat:
		format, ok := value.(string)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[types.TimeFormat] = format
	case types.ProcessingMode:
		mode, ok := value.(string)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[types.ProcessingMode] = mode
	case types.LoggerHook:
		enable, ok := value.(bool)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[types.LoggerHook] = enable
	case types.LoggerInstance:
		logger, ok := value.(*logrus.Logger)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[types.LoggerHook] = logger
	case types.ValidatorEnable:
		enable, ok := value.(bool)
		if !ok {
			return ErrInvalidOptions
		}
		c.options[types.ValidatorEnable] = enable
	default:
		return ErrInvalidOptionsType
	}

	return nil

}

func (c *Configurator) setDefaultOptions() {
	// Time format
	c.options[types.TimeFormat] = time.RFC3339
	// Processing mode
	c.options[types.ProcessingMode] = types.OverwritingMode
	// Logger hook
	c.options[types.LoggerHook] = false
	// Logger
	c.options[types.LoggerInstance] = logrus.StandardLogger()
	// Validator
	c.options[types.ValidatorEnable] = true
}

func (c *Configurator) getTimeFormat() string {
	return c.options[types.TimeFormat].(string)
}

func (c *Configurator) getProcessingMode() string {
	return c.options[types.ProcessingMode].(string)
}

func (c *Configurator) getLoggerHook() bool {
	return c.options[types.LoggerHook].(bool)
}

func (c *Configurator) getLogger() *logrus.Logger {
	return c.options[types.LoggerInstance].(*logrus.Logger)
}

func (c *Configurator) getValidatorEnable() bool {
	return c.options[types.ValidatorEnable].(bool)
}
