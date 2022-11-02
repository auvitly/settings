//go:build !windows
// +build !windows

package internal

import (
	"log/syslog"
	"settings/types"

	"github.com/sirupsen/logrus"
	logrusSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

// syslogHook создает хук для сислога
func syslogHook(loggerSettings types.Logger) (logrus.Hook, error) {

	switch loggerSettings.SyslogProtocol {
	case "udp", "tcp":
		return logrusSyslog.NewSyslogHook(
			loggerSettings.SyslogProtocol,
			loggerSettings.Syslog,
			syslog.Priority(loggerSettings.SysLogLevel),
			"")
	case "loc":
		return logrusSyslog.NewSyslogHook("", "", syslog.Priority(loggerSettings.SysLogLevel), "")
	default:
		return nil, nil
	}
}
