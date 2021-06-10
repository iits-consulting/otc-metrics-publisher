package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func logSettings() {
	// default
	zerolog.TimestampFieldName = "timestamp"

	// set up logs
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// JSON or CBOR output
	if strings.EqualFold(cfg.LogFormat, "plain") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	// log level
	switch {
	case strings.EqualFold(cfg.LogLevel, "trace") == true:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case strings.EqualFold(cfg.LogLevel, "debug") == true:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case strings.EqualFold(cfg.LogLevel, "warn") == true:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case strings.EqualFold(cfg.LogLevel, "warning") == true:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case strings.EqualFold(cfg.LogLevel, "info") == true:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case strings.EqualFold(cfg.LogLevel, "error") == true:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case strings.EqualFold(cfg.LogLevel, "fatal") == true:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case strings.EqualFold(cfg.LogLevel, "panic") == true:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case strings.EqualFold(cfg.LogLevel, "quite") == true:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case strings.EqualFold(cfg.LogLevel, "nolog") == true:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Trace().Msgf("%#v", cfg)
}
