package main

import (
	"errors"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"

	"github.com/rs/zerolog/log"
)

func sendMetrics() (err error) {
	if len(metrics) == 0 {
		err = errors.New("metrics list is empty")
		return
	}

	opts := golangsdk.AKSKAuthOptions{
		IdentityEndpoint: cfg.AuthEndpoint,
		AccessKey:        cfg.AccessKey,
		SecretKey:        cfg.SecretKey,
		ProjectId:        cfg.ProjectID,
	}

	client, err := openstack.AuthenticatedClient(opts)

	sc, err := openstack.NewCESClient(client, golangsdk.EndpointOpts{})
	if err != nil {
		log.Err(err).Msg("openstack new ces client failed")
		return
	}

	m := metricdata.AddMetricDataOpts(metrics)

	log.Debug().Msgf("send metrics: %+v", m)

	res := metricdata.AddMetricData(sc, m)
	log.Trace().Msgf("res: %+v", res)

	if res.Err != nil {
		log.Trace().Msgf("send body: %+v", res.Body)
		err = res.Err
		return
	}

	lastSend = time.Now()
	return
}

func envelopMetric(name string, value float64, unit string) (m metricdata.AddMetricDataItem) {
	var d metricdata.MetricsDimension
	if len(cfg.InstanceID) > 0 {
		d.Name = "instance_id"
		d.Value = cfg.InstanceID
	}
	m.Metric.Dimensions = append(m.Metric.Dimensions, d)

	m.Metric.Namespace = cfg.NameSpace
	m.Metric.MetricName = name
	m.Ttl = cfg.TTL
	m.CollectTime = int(time.Now().UTC().UnixNano() / int64(time.Millisecond))
	m.Type = "float"
	m.Unit = unit
	m.Value = value

	log.Trace().Msgf("%#v", m)

	return
}

func cleanMetrics() {
	mutex.Lock()
	for i := len(metrics) - 1; i >= 0; i-- {
		if metrics[i].CollectTime < int(lastSend.UTC().UnixNano()/int64(time.Millisecond)) {
			metrics = append(metrics[:i], metrics[i+1:]...)
		}
	}
	mutex.Unlock()
}
