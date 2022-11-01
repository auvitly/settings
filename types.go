package config

import "settings/internal"

const (
	TimeFormat     Options = "time_format"
	ProcessingMode Options = "processing_mode"
	LoggerHook     Options = "logger_hook"
	LoggerInstance Options = "logger_instance"
)

const (
	OverwritingMode = internal.OverwritingMode
	ComplementMode  = internal.ComplementMode
)

type Options string

type Logger struct {
	internal.Logger
}

type SyslogLevel internal.SyslogLevel
