//go:build windows
// +build windows

package internal

import (
	"github.com/sirupsen/logrus"
	"settings/types"
)

type SyslogLevel int

// syslogHook создает пустой хук для сислога (на самом деле нет)
func syslogHook(loggerSettings types.Logger) (hook logrus.Hook, err error) {
	return
}

// parseSyslog заглушает обработку сислога под Windows
func (h *Handler) parseLogger() (int64, error) {
	return 0, nil
}
