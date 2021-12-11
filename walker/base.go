package walker

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"s3-exporter/stats"
	"s3-exporter/utils"
	"strings"
)

type BaseWalkerConfig struct {
	Depth              uint              `long:"maxDepth" required:"true" default:"1" env:"MAX_DEPTH" description:"Maximum lookup depth; Will be used to group paths and results"`
	BinNumber          int               `long:"histogram-bins" required:"true" default:"30" env:"HISTOGRAM_BINS" description:"Number of bins for histograms"`
	BinStart           float64           `long:"histogram-start" required:"true" default:"10_000_000" env:"HISTOGRAM_START" description:"Value of first bin in bytes"`
	BinIncrementFactor float64           `long:"histogram-factor" required:"true" default:"1.5" env:"HISTOGRAM_FACTOR" description:"How much do we increase the size of bins (exponentially)"`
	PrefixFilters      []string          `long:"prefix-filter" required:"false" env:"PREFIX_FILTER" description:"Prefixes or part of prefix to be ignored"`
	CustomLabels       map[string]string `long:"custom-labels" env:"CUSTOM_LABELS" description:"Labels to add for prometheus exporters"`
}

type baseWalker struct {
	config        *BaseWalkerConfig
	Stats         stats.StatsInterface
	blockFlag     bool
	prefixPattern []*regexp.Regexp
}

func (b *baseWalker) Init(config Config, labels map[string]string, labelsNames []string) error {
	err := b.ValidateConfig(config)
	if err != nil {
		return err
	}
	b.config = config.BaseWalkerConfig

	b.prefixPattern = utils.BuildPatternsFromStrings(b.config.PrefixFilters)
	b.Stats = stats.NewPrometheusStatsHolder(labels, labelsNames, b.config.BinStart, b.config.BinIncrementFactor, b.config.BinNumber)
	b.blockFlag = false
	return nil
}

func (b *baseWalker) ValidateConfig(config Config) error {
	return nil
}

func (b *baseWalker) Walk() error {
	err := fmt.Errorf("impossible to call the abstract base method directly")
	log.Panicln(err)
	return err
}

func (b *baseWalker) ProcessFile(base string, path string, size int64, depth uint, contentType string, labels map[string]string) {

	log.Tracef("Current file %s", path)
	nobase := strings.TrimPrefix(path, base)
	fp := strings.TrimPrefix(filepath.ToSlash(nobase), filepath.VolumeName(path))
	currentDepth := strings.Split(fp, "/")

	usableDepth := uint(len(currentDepth) - 1)
	if usableDepth >= (depth + 1) {
		usableDepth = depth + 1
	}

	var prefix string
	if len(currentDepth) == 1 {
		prefix = "ROOT"
	} else {
		prefix = strings.Join(currentDepth[0:usableDepth], "/")
	}
	log.Debug("Path: ", fp, " Size :", size, " Prefix :", prefix)

	if utils.MatchExclude(b.prefixPattern, prefix) {
		log.Debug("Excluded ", fp)
		return
	}

	b.Stats.ProcessFile(prefix, uint64(size), uint64(len(currentDepth)), filepath.Ext(path), contentType, labels)
}

/*
getMimeType is Too Slow !
*/
func (b *baseWalker) getMimeType(path string) string {

	f, err := os.Open(path)
	if err != nil {
		log.Error("Could not open file: ", err)
		return ""
	}
	defer f.Close()

	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		log.Error("Could not read 512 bytes: ", err)
		return ""
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer)
}

func (b *baseWalker) startProcessing() {
	b.Stats.StartProcessing()
}

func (b *baseWalker) endProcessing() {
	b.Stats.EndProcessing()
}
