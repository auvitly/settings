//go:build !windows
// +build !windows

package types

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log/syslog"
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

type SyslogLevel syslog.Priority

// ParseSyslogPriority конвертирует уровень логирования для syslog.
func ParseSyslogPriority(lvl string) (SyslogLevel, error) {

	switch strings.ToLower(lvl) {
	case "panic":
		return SyslogLevel(syslog.LOG_EMERG), nil
	case "fatal":
		return SyslogLevel(syslog.LOG_CRIT), nil
	case "error":
		return SyslogLevel(syslog.LOG_ERR), nil
	case "warn", "warning":
		return SyslogLevel(syslog.LOG_WARNING), nil
	case "info":
		return SyslogLevel(syslog.LOG_INFO), nil
	case "debug":
		return SyslogLevel(syslog.LOG_DEBUG), nil
	case "trace":
		return SyslogLevel(syslog.LOG_NOTICE), nil
	}

	return 0, errors.Errorf("unknown syslog level: %s", lvl)
}
