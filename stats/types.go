package stats

type StatsInterface interface {
	ProcessFile(prefix string, size uint64, depth uint64, ext string, contentType string, labels map[string]string)
	EndProcessing()
	StartProcessing()
	Reset()
}
