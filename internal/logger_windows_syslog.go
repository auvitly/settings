//go:build windows
// +build windows

package internal

import (
	"github.com/Auvitly/settings/types"
	"github.com/sirupsen/logrus"
)

// syslogHook создает пустой хук для сислога (на самом деле нет)
func syslogHook(loggerSettings types.Logger) (hook logrus.Hook, err error) {
	return
}
