package config

import (
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"s3-exporter/walker"
	"time"
)

type Config struct {
	WalkerType     string              `long:"type" description:"Walker type" env:"WALKER_TYPE" required:"true" choice:"s3" choice:"fs"`
	Walker         walker.Config       `group:"Walkers configuration" namespace:"walker" env-namespace:"WALKER"`
	Server         ServerConfiguration `group:"HTTP Server configuration" namespace:"http" env-namespace+:"HTTP"`
	ScrapeInterval time.Duration       `long:"interval" default:"10m" env:"SCRAPE_INTERVAL" required:"false" description:"Define the minimum delay between scrapes. Set this to a reasonable value to avoid unnecessary stress on drives"`
	LogLevel       string              `long:"logLevel" default:"debug" env:"LOG_LEVEL" required:"false" description:"Level for logger; available options are: debug, info, warning, error" `
}

type ServerConfiguration struct {
	Port              int    `long:"port" default:"6535" env:"PORT" required:"true" description:"HTTP(s) server port"`
	Listen            string `long:"addr" default:"" env:"ADDR" required:"true" description:"HTTP(s) listen address"`
	ServerKeyFile     string `long:"keyFile" required:"false" env:"KEY_FILE" description:"Required along with certFile to enable HTTPS"`
	ServerCertificate string `long:"certFile" required:"false" env:"CERT_FILE" description:"Required along with keyFile to enable HTTPS"`
}

func LoadConfig() *Config {
	var opts Config
	log.Info("Loading configuration from arg or environment variables")
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(1)
		}
		log.Fatal("Could not read configuration: ", err.Error())
	}
	return &opts
}
