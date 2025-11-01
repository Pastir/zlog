package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*buildCfg)

// WithAccessFile configures access log file rotation
func WithAccessFile(path string, maxSizeMB, maxBackups, maxAgeDays int, compress bool) Option {
	return func(c *buildCfg) {
		// Validate and normalize parameters
		if maxSizeMB < 0 {
			maxSizeMB = 0
		}
		if maxBackups < 0 {
			maxBackups = 0
		}
		if maxAgeDays < 0 {
			maxAgeDays = 0
		}
		c.access = rotateCfg{
			Path:       path,
			MaxSizeMB:  maxSizeMB,
			MaxBackups: maxBackups,
			MaxAgeDays: maxAgeDays,
			Compress:   compress,
		}
	}
}

// WithErrorFile configures error log file rotation
func WithErrorFile(path string, maxSizeMB, maxBackups, maxAgeDays int, compress bool) Option {
	return func(c *buildCfg) {
		// Validate and normalize parameters
		if maxSizeMB < 0 {
			maxSizeMB = 0
		}
		if maxBackups < 0 {
			maxBackups = 0
		}
		if maxAgeDays < 0 {
			maxAgeDays = 0
		}
		c.error = rotateCfg{
			Path:       path,
			MaxSizeMB:  maxSizeMB,
			MaxBackups: maxBackups,
			MaxAgeDays: maxAgeDays,
			Compress:   compress,
		}
	}
}

// WithConsoleForAccess enables/disables console stdout output for access logs
func WithConsoleForAccess(enable bool) Option {
	return func(c *buildCfg) { c.consoleStdout = enable }
}

// WithConsoleForError enables/disables console stderr output for error logs
func WithConsoleForError(enable bool) Option {
	return func(c *buildCfg) { c.consoleStderr = enable }
}

// WithInitialLevels sets initial logging levels for access and error loggers
func WithInitialLevels(access, err zapcore.Level) Option {
	return func(c *buildCfg) {
		c.initialAccessLevel = access
		c.initialErrorLevel = err
	}
}

// WithEncoder sets custom encoder configuration
func WithEncoder(enc zapcore.EncoderConfig) Option {
	return func(c *buildCfg) { c.enc = enc }
}

// WithZapOptions sets native zap.Option for loggers
func WithZapOptions(opts ...zap.Option) Option {
	return func(c *buildCfg) {
		c.zapOpts = append(c.zapOpts, opts...)
	}
}
