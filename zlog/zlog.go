package zlog

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	// Pair represents a pair of loggers: Access for access logs and Error for error logs.
	// AccessLevel and ErrorLevel can be changed at runtime.
	Pair struct {
		Access *zap.Logger
		Error  *zap.Logger

		// AccessLevel and ErrorLevel are public and can be changed at runtime
		AccessLevel zap.AtomicLevel
		ErrorLevel  zap.AtomicLevel
	}

	rotateCfg struct {
		Path       string
		MaxSizeMB  int
		MaxBackups int
		MaxAgeDays int
		Compress   bool
	}

	buildCfg struct {
		access rotateCfg
		error  rotateCfg

		consoleStdout bool
		consoleStderr bool

		enc     zapcore.EncoderConfig
		zapOpts []zap.Option

		initialAccessLevel zapcore.Level
		initialErrorLevel  zapcore.Level
	}
)

// Sync flushes any buffered log entries. Applications should take care to call Sync before exiting.
func (p *Pair) Sync() error {
	var errs []error
	if p.Access != nil {
		if err := p.Access.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	if p.Error != nil {
		if err := p.Error.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// Use fmt.Errorf or create a simple error message
		return &syncError{errs: errs}
	}
	return nil
}

type syncError struct {
	errs []error
}

func (e *syncError) Error() string {
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}
	return "multiple sync errors occurred"
}

func (e *syncError) Unwrap() []error {
	return e.errs
}

func defaultEncoder() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newRotateWriter(c rotateCfg) zapcore.WriteSyncer {
	if c.Path == "" {
		// Empty path means discard logs
		return zapcore.AddSync(io.Discard)
	}
	// lumberjack MaxSize is in megabytes
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   c.Path,
		MaxSize:    c.MaxSizeMB,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAgeDays,
		Compress:   c.Compress,
	})
}

func makeCore(encCfg zapcore.EncoderConfig, ws zapcore.WriteSyncer, lvl zap.AtomicLevel) zapcore.Core {
	return zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), ws, lvl)
}

func tee(ws1, ws2 zapcore.WriteSyncer) zapcore.WriteSyncer {
	switch {
	case ws1 == nil:
		return ws2
	case ws2 == nil:
		return ws1
	default:
		return zapcore.NewMultiWriteSyncer(ws1, ws2)
	}
}

// New returns a pair of loggers (access/error)
func New(opts ...Option) (*Pair, error) {
	cfg := buildCfg{
		access:             rotateCfg{},
		error:              rotateCfg{},
		consoleStdout:      false,
		consoleStderr:      false,
		enc:                defaultEncoder(),
		initialAccessLevel: zapcore.InfoLevel,
		initialErrorLevel:  zapcore.ErrorLevel,
		zapOpts:            []zap.Option{},
	}
	for _, o := range opts {
		o(&cfg)
	}

	// levels
	accessLevel := zap.NewAtomicLevelAt(cfg.initialAccessLevel)
	errorLevel := zap.NewAtomicLevelAt(cfg.initialErrorLevel)

	// writers
	accessFile := newRotateWriter(cfg.access)
	errorFile := newRotateWriter(cfg.error)

	var accessConsole zapcore.WriteSyncer
	if cfg.consoleStdout {
		accessConsole = zapcore.AddSync(os.Stdout)
	}
	var errorConsole zapcore.WriteSyncer
	if cfg.consoleStderr {
		errorConsole = zapcore.AddSync(os.Stderr)
	}

	// cores (tee: file + console)
	accessCore := makeCore(cfg.enc, tee(accessFile, accessConsole), accessLevel)
	errorCore := makeCore(cfg.enc, tee(errorFile, errorConsole), errorLevel)

	errOpts := append([]zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}, cfg.zapOpts...)

	access := zap.New(accessCore, cfg.zapOpts...)
	errorL := zap.New(errorCore, errOpts...)

	return &Pair{
		Access:      access,
		Error:       errorL,
		AccessLevel: accessLevel,
		ErrorLevel:  errorLevel,
	}, nil
}
