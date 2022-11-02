//go:build windows
// +build windows

package types

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
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

type SyslogLevel int

// ParseSyslogPriority конвертирует уровень логирования для syslog.
func ParseSyslogPriority(lvl string) (SyslogLevel, error) {

	switch strings.ToLower(lvl) {
	case "panic":
		return SyslogLevel(0), nil
	case "fatal":
		return SyslogLevel(1), nil
	case "error":
		return SyslogLevel(2), nil
	case "warn", "warning":
		return SyslogLevel(3), nil
	case "info":
		return SyslogLevel(4), nil
	case "debug":
		return SyslogLevel(5), nil
	case "trace":
		return SyslogLevel(6), nil
	}

	return 0, errors.Errorf("unknown syslog level: %s", lvl)
}
