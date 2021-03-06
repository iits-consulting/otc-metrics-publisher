package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	"github.com/rs/zerolog/log"
)

var (
	cfg      config
	metrics  []metricdata.AddMetricDataItem
	lastSend time.Time
	mutex    = &sync.Mutex{}

	// parameters set by CI pipeline
	version = ""
	date    = ""
	commit  = ""
)

func main() {
	// command line options
	debug := flag.Bool("debug", false, "debug mode")
	trace := flag.Bool("trace", false, "trace mode")
	v := flag.Bool("version", false, "print version")
	flag.Parse()

	if *v == true {
		// version request
		fmt.Printf("version: %s, commit: %s, date: %s\n", version, commit, date)
		os.Exit(0)
	}

	err := envconfig.Process(envVarPrefix, &cfg)
	if err != nil {
		log.Err(err).Msg("error reading config values")
		os.Exit(1)
	}

	// overwrite log level
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

	// start goroutine to grab metrics
	go func() {
		var listOfScripts []fs.FileInfo
		var err error

		// get list of files
		// that will generate non-changeble list of scrpts on the start
		if cfg.RefreshScriptsList == false {
			listOfScripts, err = ioutil.ReadDir(cfg.ScriptsDir)
			if err != nil {

			}
		}

		for {
			// get list of files
			// that allows to add and remove scripts ad-hoc without daemon restart
			if cfg.RefreshScriptsList == true {
				listOfScripts, err = ioutil.ReadDir(cfg.ScriptsDir)
			}

			if err == nil {
				for _, file := range listOfScripts {
					// if file is not dir, bigger than 0 and is executable
					if file.IsDir() == false &&
						file.Size() > 0 &&
						file.Mode()&0111 == 0111 {
						go scriptExec(file.Name())
						time.Sleep(time.Millisecond * 10) // sleep 10ms before launching next script
					}
				}
			} else {
				log.Err(err).Msgf("cannot fetch list of scripts in %s directory", cfg.ScriptsDir)
			}

			time.Sleep(time.Second * time.Duration(cfg.GrabInterval))

		}
	}()

	// endless loop to send and clean the metrics up
	for {
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
