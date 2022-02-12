package stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/willena/s3-exporter/utils"
	"time"
)

const METRICS_GROUP = "file_walker"

type PrometheusStats struct {
	// Simple stats
	MaxDepth          *prometheus.GaugeVec
	CollectDuration   *prometheus.GaugeVec
	TotalObjectsSize  *prometheus.GaugeVec
	TotalObjectsCount *prometheus.GaugeVec
	LastWalkStart     *prometheus.GaugeVec

	//Per prefix stats
	PerPrefixObjectsSizeHistogram      *prometheus.HistogramVec
	PerPrefixObjectsSize               *prometheus.GaugeVec
	PerPrefixObjectsCount              *prometheus.GaugeVec
	PerPrefixPerExtensionObjectCount   *prometheus.GaugeVec
	PerPrefixPerExtensionObjectsSize   *prometheus.GaugeVec
	PerPrefixPerContentTypeObjectCount *prometheus.GaugeVec
	PerPrefixPerContentTypeObjectsSize *prometheus.GaugeVec

	startTime time.Time
	maxDepth  uint64

	currentMetrics                []prometheus.Collector
	constLabels                   prometheus.Labels
	namesWithPrefix               []string
	namesWithPrefixAndContentType []string
	start                         float64
	namesWithPrefixAndExt         []string
	factor                        float64
	number                        int
	names                         []string
}

func (p *PrometheusStats) ProcessFile(prefix string, size uint64, depth uint64, ext string, contentType string, labels map[string]string) {

	if depth >= p.maxDepth {
		p.maxDepth = depth
		p.MaxDepth.With(labels).Set(float64(depth))
	}

	p.TotalObjectsCount.With(labels).Add(1)
	p.TotalObjectsSize.With(labels).Add(float64(size))

	prefixLabel := utils.MergeMapsRight(prometheus.Labels{
		"prefix": prefix,
	}, labels)

	p.PerPrefixObjectsSizeHistogram.With(prefixLabel).Observe(float64(size))
	p.PerPrefixObjectsSize.With(prefixLabel).Add(float64(size))
	p.PerPrefixObjectsCount.With(prefixLabel).Add(1)

	prefixExtLabels := utils.MergeMapsRight(prometheus.Labels{
		"prefix": prefix,
		"ext":    ext,
	}, labels)

	p.PerPrefixPerExtensionObjectCount.With(prefixExtLabels).Add(1)
	p.PerPrefixPerExtensionObjectsSize.With(prefixExtLabels).Add(float64(size))

	prefixContentTypeLabels := utils.MergeMapsRight(prometheus.Labels{
		"prefix":      prefix,
		"contentType": contentType,
	}, labels)

	p.PerPrefixPerContentTypeObjectCount.With(prefixContentTypeLabels).Add(1)
	p.PerPrefixPerContentTypeObjectsSize.With(prefixContentTypeLabels).Add(float64(size))

}

func (p *PrometheusStats) EndProcessing() {
	p.CollectDuration.With(nil).Set(float64(time.Now().Sub(p.startTime)))
	p.updatePrometheusGauges()
}

func (p *PrometheusStats) StartProcessing() {
	p.Reset()
	p.startTime = time.Now()
	p.LastWalkStart.With(nil).Set(float64(p.startTime.Unix()))
}

func (p *PrometheusStats) Reset() {
	// Create new Gauges each time...
	// Old ones should still be available during refresh.
	p.CollectDuration = createGaugeVect("stats_collection_duration", "Time spent reading object and folders", p.constLabels, nil)
	p.LastWalkStart = createGaugeVect("stats_collection_date", "Date when the stats collection started", p.constLabels, nil)
	p.MaxDepth = createGaugeVect("max_tree_depth", "Maximum depth of folder tree", p.constLabels, p.names)
	p.TotalObjectsSize = createGaugeVect("total_objects_size", "Total objects volume in bytes", p.constLabels, p.names)
	p.TotalObjectsCount = createGaugeVect("total_objects_count", "total number of objects found", p.constLabels, p.names)
	p.PerPrefixObjectsSizeHistogram = createHistogramVect("objects_sizes_count", "Histogram showing the files size repartition across prefixes", p.constLabels, p.start, p.factor, p.number, p.namesWithPrefix)
	p.PerPrefixObjectsSize = createGaugeVect("objects_size", "Objects volume across prefixes", p.constLabels, p.namesWithPrefix)
	p.PerPrefixObjectsCount = createGaugeVect("objects_count", "Objects count across prefixes", p.constLabels, p.namesWithPrefix)
	p.PerPrefixPerExtensionObjectCount = createGaugeVect("objects_extensions_count", "Repartition of objects per file extension", p.constLabels, p.namesWithPrefixAndExt)
	p.PerPrefixPerExtensionObjectsSize = createGaugeVect("objects_extensions_size", "Total size of objects per extension", p.constLabels, p.namesWithPrefixAndExt)
	p.PerPrefixPerContentTypeObjectCount = createGaugeVect("objects_content_type_count", "Repartition of objects per file ContentType", p.constLabels, p.namesWithPrefixAndContentType)
	p.PerPrefixPerContentTypeObjectsSize = createGaugeVect("objects_content_type_size", "Total size of objects per ContentType", p.constLabels, p.namesWithPrefixAndContentType)

}

// Describe implements the prometheus.Collector interface
func (p *PrometheusStats) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range p.currentMetrics {
		metric.Describe(ch)
	}
}

// Collect implements the prometheus.Collector interface
func (p *PrometheusStats) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range p.currentMetrics {
		metric.Collect(ch)
	}

}

func (p *PrometheusStats) updatePrometheusGauges() {
	p.currentMetrics = []prometheus.Collector{}
	p.currentMetrics = append(
		p.currentMetrics,
		p.MaxDepth,
		p.CollectDuration,
		p.TotalObjectsSize,
		p.TotalObjectsCount,
		p.LastWalkStart,
		p.PerPrefixObjectsSizeHistogram,
		p.PerPrefixObjectsSize,
		p.PerPrefixObjectsCount,
		p.PerPrefixPerExtensionObjectCount,
		p.PerPrefixPerExtensionObjectsSize,
		p.PerPrefixPerContentTypeObjectCount,
		p.PerPrefixPerContentTypeObjectsSize,
	)
}

func createGaugeVect(name string, help string, labels prometheus.Labels, names []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        METRICS_GROUP + "_" + name,
		Help:        help,
		ConstLabels: labels,
	}, names)
}

func createHistogramVect(name, help string, labels prometheus.Labels, start, factor float64, number int, names []string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        METRICS_GROUP + "_" + name,
		Help:        help,
		ConstLabels: labels,
		Buckets:     prometheus.ExponentialBuckets(start, factor, number),
	}, names)
}

func NewPrometheusStatsHolder(constLabels prometheus.Labels, names []string, start, factor float64, number int) StatsInterface {

	namesWithPrefix := []string{"prefix"}
	namesWithPrefixAndExt := []string{"ext", "prefix"}
	namesWithPrefixAndContentType := []string{"prefix", "contentType"}

	namesWithPrefix = append(namesWithPrefix, names...)
	namesWithPrefixAndExt = append(namesWithPrefixAndExt, names...)
	namesWithPrefixAndContentType = append(namesWithPrefixAndContentType, names...)

	ps := &PrometheusStats{
		constLabels:                   constLabels,
		namesWithPrefix:               namesWithPrefix,
		namesWithPrefixAndExt:         namesWithPrefixAndExt,
		namesWithPrefixAndContentType: namesWithPrefixAndContentType,
		start:                         start,
		factor:                        factor,
		number:                        number,
		names:                         names,
	}

	return ps
}
