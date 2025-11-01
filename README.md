# zlog

A structured logging package for Go built on top of [zap](https://github.com/uber-go/zap). `zlog` provides a dual-logger system with separate access and error loggers, file rotation, and flexible configuration options.

## Features

- **Dual Logger System**: Separate loggers for access logs and error logs
- **File Rotation**: Automatic log file rotation using [lumberjack](https://github.com/natefinch/lumberjack) with configurable size, backup, and retention policies
- **Console Output**: Optional console output (stdout for access logs, stderr for error logs)
- **Runtime Log Levels**: Dynamically adjustable log levels via atomic level controls
- **JSON Encoding**: Structured JSON logging by default with customizable encoder configuration
- **Zap Integration**: Full compatibility with zap's native options and features

## Installation

```bash
go get github.com/Pastir/zlog
```

## Quick Start

```go
package main

import (
    "log"
    
    "go.uber.org/zap"
    "github.com/Pastir/zlog"
)

func main() {
    // Create a logger pair with file rotation
    pair, err := zlog.New(
        zlog.WithAccessFile("/var/log/app/access.log", 100, 5, 30, true),
        zlog.WithErrorFile("/var/log/app/error.log", 100, 10, 30, true),
        zlog.WithConsoleForAccess(true),
        zlog.WithConsoleForError(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer pair.Sync()

    // Use the loggers
    pair.Access.Info("User logged in", zap.String("user_id", "12345"))
    pair.Error.Error("Database connection failed", zap.String("db", "main"))
}
```

## Configuration Options

### File Rotation

Configure file rotation for access and error logs separately:

```go
zlog.WithAccessFile(path string, maxSizeMB, maxBackups, maxAgeDays int, compress bool)
zlog.WithErrorFile(path string, maxSizeMB, maxBackups, maxAgeDays int, compress bool)
```

Parameters:
- `path`: Log file path (empty string disables file logging)
- `maxSizeMB`: Maximum size in megabytes before rotation (0 = no limit)
- `maxBackups`: Maximum number of backup files to keep (0 = keep all)
- `maxAgeDays`: Maximum age in days before deleting old logs (0 = no age limit)
- `compress`: Whether to compress rotated log files

### Console Output

Enable console output for debugging or development:

```go
zlog.WithConsoleForAccess(enable bool)  // stdout
zlog.WithConsoleForError(enable bool)   // stderr
```

### Log Levels

Set initial log levels for both loggers:

```go
import "go.uber.org/zap/zapcore"

zlog.WithInitialLevels(zapcore.InfoLevel, zapcore.ErrorLevel)
```

Log levels can be changed at runtime:

```go
pair.AccessLevel.SetLevel(zapcore.DebugLevel)
pair.ErrorLevel.SetLevel(zapcore.WarnLevel)
```

### Encoder Configuration

Customize the JSON encoder:

```go
import "go.uber.org/zap/zapcore"

customEncoder := zapcore.EncoderConfig{
    TimeKey:        "timestamp",
    LevelKey:       "level",
    MessageKey:     "message",
    EncodeLevel:    zapcore.LowercaseLevelEncoder,
    EncodeTime:     zapcore.RFC3339TimeEncoder,
}

zlog.WithEncoder(customEncoder)
```

### Zap Options

Add native zap options:

```go
import "go.uber.org/zap"

zlog.WithZapOptions(
    zap.AddCallerSkip(1),
    zap.Development(),
)
```

## Default Behavior

- **Access Logger**: Info level, no file rotation by default
- **Error Logger**: Error level, includes caller and stacktrace, no file rotation by default
- **Encoding**: JSON with ISO8601 timestamps, capital level names, and short caller info
- **Console**: Disabled by default

## Type Reference

### Pair

The main logger pair structure:

```go
type Pair struct {
    Access      *zap.Logger      // Access logger instance
    Error       *zap.Logger      // Error logger instance
    AccessLevel zap.AtomicLevel  // Runtime-adjustable access log level
    ErrorLevel  zap.AtomicLevel  // Runtime-adjustable error log level
}
```

### Methods

- `Sync() error`: Flushes any buffered log entries. Should be called before application exit.

## Examples

### Development Setup (Console Only)

```go
pair, err := zlog.New(
    zlog.WithConsoleForAccess(true),
    zlog.WithConsoleForError(true),
    zlog.WithInitialLevels(zapcore.DebugLevel, zapcore.DebugLevel),
)
```

### Production Setup (File Rotation Only)

```go
pair, err := zlog.New(
    zlog.WithAccessFile("/var/log/app/access.log", 500, 10, 90, true),
    zlog.WithErrorFile("/var/log/app/error.log", 500, 20, 90, true),
    zlog.WithInitialLevels(zapcore.InfoLevel, zapcore.ErrorLevel),
)
```

### Hybrid Setup (Files + Console)

```go
pair, err := zlog.New(
    zlog.WithAccessFile("/var/log/app/access.log", 100, 5, 30, true),
    zlog.WithErrorFile("/var/log/app/error.log", 100, 10, 30, true),
    zlog.WithConsoleForAccess(true),
    zlog.WithConsoleForError(true),
)
```

## Dependencies

- [go.uber.org/zap](https://github.com/uber-go/zap) - High-performance structured logging
- [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack) - Log file rotation

## License

This package uses the same license as its dependencies. Please refer to the individual dependency licenses for details.

