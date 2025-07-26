package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
)

const (
	debugLevel = "DEBUG"
	infoLevel  = "INFO"
	warnLevel  = "WARN"
	errorLevel = "ERROR"
)

func InitLogger(cfg *config.Config) *slog.Logger {
	var level slog.Leveler
	source := false

	switch cfg.Logger {
	case debugLevel:
		level = slog.LevelDebug
		source = true

	case infoLevel:
		level = slog.LevelInfo

	case warnLevel:
		level = slog.LevelWarn

	case errorLevel:
		level = slog.LevelError

	default:
		log.Fatal(errs.ErrLoggerLevel())
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: source,
	}))
}
