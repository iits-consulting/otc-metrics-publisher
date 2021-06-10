package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logLevels = map[string]zerolog.Level{
	"debug":   zerolog.DebugLevel,
	"info":    zerolog.InfoLevel,
	"warn":    zerolog.WarnLevel,
	"warning": zerolog.WarnLevel,
	"error":   zerolog.ErrorLevel,
	"fatal":   zerolog.FatalLevel,
	"panic":   zerolog.PanicLevel,
	"quite":   zerolog.Disabled,
	"nolog":   zerolog.Disabled,
	"trace":   zerolog.TraceLevel,
}

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
	if val, ok := logLevels[strings.ToLower(cfg.LogLevel)]; ok {
		zerolog.SetGlobalLevel(val)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// trace print config content
	log.Trace().Msgf("%#v", cfg)
}
