//go:build !windows && !plan9
// +build !windows,!plan9

package internal

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	logrusSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/spf13/viper"
)

type SyslogLevel = syslog.Priority

// syslogHook создает хук для сислога
func syslogHook(loggerSettings Logger) (logrus.Hook, error) {

	switch loggerSettings.SyslogProtocol {
	case "udp", "tcp":
		return logrusSyslog.NewSyslogHook(
			loggerSettings.SyslogProtocol,
			loggerSettings.Syslog,
			loggerSettings.SysLogLevel,
			"")
	case "loc":
		return logrusSyslog.NewSyslogHook("", "", loggerSettings.SysLogLevel, "")
	default:
		return nil, nil
	}
}

// parseSyslog возвращает максимальное значение типов uint.
func (h *Handler) parseSyslog(v *viper.Viper) (int64, error) {

	var (
		priority syslog.Priority
		err      error
	)

	if h.hasEnvTag {
		priority, err = ParseSyslogPriority(h.value)
		if err != nil {
			return 0, err
		}
	} else {
		priority, err = ParseSyslogPriority(v.GetString(h.tomlPath))
		if err != nil {
			return 0, err
		}
	}

	return int64(priority), nil
}
