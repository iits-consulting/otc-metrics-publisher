package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	"github.com/rs/zerolog/log"
)

func scriptExec(script string) {
	// exec timeout: https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ScriptExecTimeout)*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, cfg.ScriptsDir+"/"+script)

	// This time we can simply use Output() to get the result.
	output, err := cmd.Output()
	log.Trace().Msgf("exec file: %s, output: %#v", script, string(output))

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		log.Err(context.DeadlineExceeded).Msgf("Script timed out: %s", script)
		return
	}

	if err != nil {
		log.Err(err).Msgf("error running script %s", script)
		return
	}

	// parse output to metrics structure
	m, errParseOut := scriptParseOutput(string(output))
	log.Trace().Msgf("parsed metric object: %#v", m)

	if errParseOut == nil {
		mutex.Lock()
		metrics = append(metrics, m)
		mutex.Unlock()
	} else {
		log.Err(errParseOut).Msg("parsing output error")
		return
	}
}

func scriptParseOutput(output string) (m metricdata.AddMetricDataItem, err error) {
	// trim string
	output = strings.TrimSpace(output)
	log.Trace().Msgf("trimmed string: %s", output)

	// check empty string
	if len(output) == 0 {
		err = errors.New("empty output")
		return
	}

	c := struct {
		Name  string  `json:"name" required:"true"`
		Value float64 `json:"value" required:"true"`
		Unit  string  `json:"unit,omitempty"`
	}{}

	// check if output is JSON
	if output[0] == "{"[0] && output[len(output)-1] == "}"[0] {
		err = json.Unmarshal([]byte(output), &c)
		if err == nil {
			m = envelopMetric(c.Name, c.Value, c.Unit)
		}
		return
	}

	// check if multiline output
	if strings.Count(output, "\n") > 1 &&
		strings.Count(output, "\n") < 4 { // magic number, should be 2-3 strings
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			akv := strings.Split(scanner.Text(), "=")
			if len(akv) == 2 { // 2 is key = value
				k := strings.TrimSpace(akv[0])
				v := strings.TrimSpace(akv[1])
				if strings.EqualFold(k, "name") {
					c.Name = v
				}
				if strings.EqualFold(k, "value") {
					c.Value, err = strconv.ParseFloat(v, 64)
					if err != nil {
						return
					}
				}
				if strings.EqualFold(k, "unit") {
					c.Unit = v
				}
			}
		}

		// name lenght after parsing multiline output
		if len(c.Name) == 0 {
			err = errors.New("wrong name in multiline output: " + output)
		}
		return
	}

	// check a single-line output
	if strings.Count(output, "\n") < 2 { // this makes sense if line is not finishing with "\n"
		// split
		arr := strings.Split(output, ";")
		if len(arr) < 2 || len(arr) > 5 {
			err = errors.New("unable to split the string on 2-4 fields:" + output)
			return
		}

		// get metric name
		if len(arr[0]) > 0 {
			c.Name = arr[0]
		} else {
			err = errors.New("metric name has 0 size")
			return
		}

		// get metric value
		c.Value, err = strconv.ParseFloat(arr[1], 64)
		if err != nil {
			return
		}

		// get unit value
		if len(arr) > 2 {
			c.Unit = arr[2]
		}
	}

	m = envelopMetric(c.Name, c.Value, c.Unit)
	return

}
