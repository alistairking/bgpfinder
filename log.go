package bgpfinder

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// TODO: this is just a sketch. do useful things here
// TODO: consider moving to internal package?

type LoggerConfig struct {
	LogLevel string `help:"Log level" default:"info"`
}

type Logger struct {
	zerolog.Logger
}

func NewLogger(cfg LoggerConfig) (*Logger, error) {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	zl := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
	l := &Logger{
		Logger: zl,
	}
	return l, nil
}

// Create a sub-logger with the given module name
func (l Logger) ModuleLogger(module string) Logger {
	l.Logger = l.With().
		Str("package", "bgpfinder").
		Str("module", module).
		Logger()
	return l
}
