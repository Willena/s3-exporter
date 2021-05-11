package utils

import log "github.com/sirupsen/logrus"

func InitLogger(level string) {
	log.SetFormatter(&log.TextFormatter{})
	var logLevel = log.DebugLevel
	if level != "" {
		level, err := log.ParseLevel(level)
		if err == nil {
			logLevel = level
		}
	}
	log.SetLevel(logLevel)
	log.Infof("Logger is now at %s level", logLevel)
}
