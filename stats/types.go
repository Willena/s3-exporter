package stats

import "github.com/prometheus/client_golang/prometheus"

type StatsInterface interface {
	ProcessFile(prefix string, size uint64, depth uint64, ext string, contentType string, labels map[string]string)
	EndProcessing()
	StartProcessing()
	Reset()
	Describe(descs chan<- *prometheus.Desc)
	Collect(metrics chan<- prometheus.Metric)
}
