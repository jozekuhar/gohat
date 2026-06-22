package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(debug bool) *slog.Logger {
	fileWriter := &lumberjack.Logger{
		Filename:   "./log/app.log",
		MaxSize:    500,
		MaxAge:     30,
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	}
	writers := []io.Writer{fileWriter}
	if debug {
		consoleWriter := os.Stdout
		writers = append(writers, consoleWriter)
	}

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if debug {
		handlerOpts.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(
		io.MultiWriter(writers...),
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	))

	return logger
}
