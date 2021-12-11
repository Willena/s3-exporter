package stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"s3-exporter/utils"
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
}

func (p *PrometheusStats) ProcessFile(prefix string, size uint64, depth uint64, ext string, contentType string, labels map[string]string) {

	if depth >= p.maxDepth {
		p.maxDepth = depth
		p.MaxDepth.With(nil).Set(float64(depth))
	}

	p.TotalObjectsCount.With(nil).Add(1)
	p.TotalObjectsSize.With(nil).Add(float64(size))

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
}

func (p *PrometheusStats) StartProcessing() {
	p.Reset()
	p.startTime = time.Now()
	p.LastWalkStart.With(nil).Set(float64(p.startTime.Unix()))
}

func (p *PrometheusStats) Reset() {
	// Simple stats
	p.maxDepth = 0
	p.MaxDepth.Reset()
	p.CollectDuration.Reset()
	p.TotalObjectsSize.Reset()
	p.TotalObjectsCount.Reset()
	p.LastWalkStart.Reset()

	//Per prefix stats
	p.PerPrefixObjectsSizeHistogram.Reset()
	p.PerPrefixObjectsSize.Reset()
	p.PerPrefixObjectsCount.Reset()
	p.PerPrefixPerExtensionObjectCount.Reset()
	p.PerPrefixPerExtensionObjectsSize.Reset()
	p.PerPrefixPerContentTypeObjectCount.Reset()
	p.PerPrefixPerContentTypeObjectsSize.Reset()
}

func createGaugeVect(name string, help string, labels prometheus.Labels, names []string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        METRICS_GROUP + "_" + name,
		Help:        help,
		ConstLabels: labels,
	}, names)
}

func createHistogramVect(name, help string, labels prometheus.Labels, start, factor float64, number int, names []string) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        METRICS_GROUP + "_" + name,
		Help:        help,
		ConstLabels: labels,
		Buckets:     prometheus.ExponentialBuckets(start, factor, number),
	}, names)
}

func NewPrometheusStatsHolder(constLabels prometheus.Labels, names []string, start, factor float64, number int) StatsInterface {

	namesWithPrefix := append(names, "prefix")
	namesWithPrefixAndExt := append(namesWithPrefix, "ext")
	namesWithPrefixAndContentType := append(namesWithPrefix, "contentType")

	return &PrometheusStats{
		MaxDepth:                           createGaugeVect("max_tree_depth", "Maximum depth of folder tree", constLabels, nil),
		CollectDuration:                    createGaugeVect("stats_collection_duration", "Time spent reading object and folders", constLabels, nil),
		TotalObjectsSize:                   createGaugeVect("total_objects_size", "Total objects volume in bytes", constLabels, nil),
		TotalObjectsCount:                  createGaugeVect("total_objects_count", "total number of objects found", constLabels, nil),
		LastWalkStart:                      createGaugeVect("stats_collection_date", "Date when the stats collection started", constLabels, nil),
		PerPrefixObjectsSizeHistogram:      createHistogramVect("objects_sizes_count", "Histogram showing the files size repartition across prefixes", constLabels, start, factor, number, namesWithPrefix),
		PerPrefixObjectsSize:               createGaugeVect("objects_size", "Objects volume across prefixes", constLabels, namesWithPrefix),
		PerPrefixObjectsCount:              createGaugeVect("objects_count", "Objects count across prefixes", constLabels, namesWithPrefix),
		PerPrefixPerExtensionObjectCount:   createGaugeVect("objects_extensions_count", "Repartition of objects per file extension", constLabels, namesWithPrefixAndExt),
		PerPrefixPerExtensionObjectsSize:   createGaugeVect("objects_extensions_size", "Total size of objects per extension", constLabels, namesWithPrefixAndExt),
		PerPrefixPerContentTypeObjectCount: createGaugeVect("objects_content_type_count", "Repartition of objects per file ContentType", constLabels, namesWithPrefixAndContentType),
		PerPrefixPerContentTypeObjectsSize: createGaugeVect("objects_content_type_size", "Total size of objects per ContentType", constLabels, namesWithPrefixAndContentType),
	}
}
