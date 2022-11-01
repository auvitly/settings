package config

import "settings/internal"

const (
	TimeFormat     = internal.TimeFormat
	ProcessingMode = internal.ProcessingMode
	LoggerHook     = internal.LoggerHook
	LoggerInstance = internal.LoggerInstance
)

const (
	OverwritingMode = internal.OverwritingMode
	ComplementMode  = internal.ComplementMode
)

type Logger struct {
	internal.Logger
}

type SyslogLevel internal.SyslogLevel
