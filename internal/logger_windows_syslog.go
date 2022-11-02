//go:build windows
// +build windows

package internal

import (
	"github.com/sirupsen/logrus"
	"settings/types"
)

// syslogHook создает пустой хук для сислога (на самом деле нет)
func syslogHook(loggerSettings types.Logger) (hook logrus.Hook, err error) {
	return
}
