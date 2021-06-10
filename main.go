package main

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	"github.com/rs/zerolog/log"
)

const (
	envVarPrefix = "metrics"
)

var (
	cfg      config
	metrics  []metricdata.AddMetricDataItem
	lastSend time.Time
	mutex    = &sync.Mutex{}
)

func main() {
	var err error

	err = envconfig.Process(envVarPrefix, &cfg)
	if err != nil {
		log.Err(err).Msg("error reading config values")
		os.Exit(1)
	}

	// debug mode from commad line
	debug := flag.Bool("debug", false, "debug mode")
	trace := flag.Bool("trace", false, "trace mode")
	flag.Parse()
	if *debug == true {
		cfg.LogLevel = "debug"
	}
	if *trace == true {
		cfg.LogLevel = "trace"
	}

	// configure logger
	logSettings()

	// config sanity checks and default values
	err = configHandler()
	if err != nil {
		log.Err(err).Msg("missing required parameter")
		os.Exit(1)
	}

	// start basic cloud eye exporter
	go metricsHandler()

	// endless
	for {
		// sleep
		time.Sleep(time.Second * time.Duration(cfg.SendInterval))

		// regullary send metrics
		err = sendMetrics()
		if err != nil {
			// metrics sending error
			log.Err(err).Msgf("failed to send metrics")
		}
		go cleanMetrics()
	}

}
