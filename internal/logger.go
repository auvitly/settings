package internal

import (
	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/sirupsen/logrus"
	"io"
)

type Logger struct {
	LogLevel       logrus.Level `env:"LOG_LEVEL" toml:"level" json:"level" xml:"level" yaml:"level" default:"debug"`
	Syslog         string       `env:"SYSLOG" toml:"syslog_addr" json:"syslog_addr" xml:"syslog_addr" yaml:"syslog_addr" default:"127.0.0.1:514" validate:"tcp_addr"`
	SyslogProtocol string       `env:"SYSLOG_PROTOCOL" toml:"syslog_protocol" json:"syslog_protocol" xml:"syslog_protocol" yaml:"syslog_protocol" default:"udp" validate:"min=3,max=3"`
	SysLogLevel    SyslogLevel  `env:"SYSLOG_LEVEL" toml:"syslog_level" json:"syslog_level" xml:"syslog_level" yaml:"syslog_level" default:"debug"`
	Colour         bool         `env:"COLOUR" toml:"colour" json:"colour" xml:"colour" yaml:"colour" default:"false"`
	StdOut         bool         `env:"STDOUT" toml:"stdout" json:"stdout" xml:"stdout" yaml:"stdout" default:"true"`
	GraylogLevel   logrus.Level `env:"GRAYLOG_LEVEL" toml:"graylog_level" json:"graylog_level" xml:"graylog_level" yaml:"graylog_level" default:"debug"`
	Graylog        string       `env:"GRAYLOG" toml:"graylog" json:"graylog" xml:"graylog" yaml:"graylog"`
}

func (c *Configurator) configureLogger(config Logger) error {

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
