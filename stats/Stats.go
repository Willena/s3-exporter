package stats

import (
	"code.cloudfoundry.org/bytefmt"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

type StatsOptions struct {
	bankSizeIncrements uint64
	bankNumber         uint
}

type HistogramItem struct {
	Name  string
	Value uint64
}

type Stats struct {
	StartTime       time.Time
	EndTime         time.Time
	MaxDepth        uint64
	Count           map[string][]HistogramItem   // a[prefix][0...N] count
	Size            map[string]uint64            // a[prefix] 65MB
	Ext             map[string]map[string]uint64 // a[prefix]["ext"] count
	ExtSize         map[string]map[string]uint64 // a[prefix]["ext"] Size
	ContentType     map[string]map[string]uint64 // a[prefix]["image/png"] count
	ContentTypeSize map[string]map[string]uint64 // a[prefix]["image/png"] Size
	TotalObjects    uint64
	TotalSize       uint64
	Options         StatsOptions
}

func (s *Stats) ExtraComputations() {
	for _, v := range s.Size {
		s.TotalSize += v
	}

	for _, v := range s.Count {
		for i := range v {
			s.TotalObjects += v[i].Value
		}
	}
}

func (s *Stats) ProcessFile(prefix string, size uint64, depth uint64, ext string, contentType string) {
	// Update Depth
	if depth > s.MaxDepth {
		s.MaxDepth = depth
	}

	// Updated Size
	if _, ok := s.Size[prefix]; ok {
		s.Size[prefix] += size
	} else {
		s.Size[prefix] = size
	}

	//Update Count
	if val, ok := s.Count[prefix]; ok {
		val[s.getBin(size)].Value += 1
	} else {
		s.Count[prefix] = s.createBins()
		log.Debugf("Created %d bins for prefix %s", len(s.Count[prefix]), prefix)
	}

	//Update Exts count
	if val, ok := s.Ext[prefix]; ok {
		val[ext] += 1
	} else {
		s.Ext[prefix] = map[string]uint64{}
		s.Ext[prefix][ext] += 1
	}

	//Update Exts Size
	if val, ok := s.ExtSize[prefix]; ok {
		val[ext] += size
	} else {
		s.ExtSize[prefix] = map[string]uint64{}
		s.ExtSize[prefix][ext] += size
	}

	//Update ContentType count
	if val, ok := s.ContentType[prefix]; ok {
		val[ext] += 1
	} else {
		s.ContentType[prefix] = map[string]uint64{}
		s.ContentType[prefix][contentType] += 1
	}

	//Update ContentTypeSize Size
	if val, ok := s.ContentTypeSize[prefix]; ok {
		val[ext] += size
	} else {
		s.ContentTypeSize[prefix] = map[string]uint64{}
		s.ContentTypeSize[prefix][contentType] += size
	}
}
func (s *Stats) computeBins() []HistogramItem {
	var items []HistogramItem
	size := s.Options.bankSizeIncrements

	for i := uint(1); i < s.Options.bankNumber; i++ {
		items = append(items, HistogramItem{
			Name:  "<" + bytefmt.ByteSize(size),
			Value: 0,
		})
		size += s.Options.bankSizeIncrements
	}

	items = append(items, HistogramItem{
		Name:  ">" + bytefmt.ByteSize(size-s.Options.bankSizeIncrements),
		Value: 0,
	})

	return items

}

func (s *Stats) createBins() []HistogramItem {
	return s.computeBins()
}

func (s *Stats) getBin(size uint64) int {
	binN := uint(math.Ceil(float64(size / s.Options.bankSizeIncrements)))
	if binN >= s.Options.bankNumber {
		binN = s.Options.bankNumber - 1
	}

	return int(binN)

}

func NewStat(numberFileSizeBank uint, sizeIncrement uint64) *Stats {
	a := &Stats{
		Options: StatsOptions{
			bankSizeIncrements: sizeIncrement,
			bankNumber:         numberFileSizeBank,
		},
		MaxDepth:        0,                              // Gauge => max depth found
		Count:           map[string][]HistogramItem{},   // Histogram labels: prefix => size histogram (numner of item given size)
		Size:            map[string]uint64{},            // Gauge labels: prefix => size per prefix
		Ext:             map[string]map[string]uint64{}, // Gauge labels: prefix, ext => count of ext
		ExtSize:         map[string]map[string]uint64{}, // Gauge labels: prefix, ext => total size of ext
		ContentType:     map[string]map[string]uint64{}, // Gauge labels: prefix, contentType => count per content type
		ContentTypeSize: map[string]map[string]uint64{}, // Gauge labels: prefix, contentType => size per content type
		// Gauge => scrape duration ms
		// Gauge => total size
		// Gauge => total objects
		// Counter => last walk
	}

	return a
}
