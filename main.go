package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/willena/s3-exporter/config"
	"github.com/willena/s3-exporter/utils"
	"github.com/willena/s3-exporter/walker"
	"net/http"
	"strconv"
	"time"
)

func main() {
	utils.InitLogger("")
	log.Info("Starting S3Exporter...")

	opts := config.LoadConfig()

	utils.InitLogger(opts.LogLevel)

	walkerInst, err := walker.FromConfig(opts.Walker, opts.WalkerType)
	if err != nil {
		log.Fatal(err.Error())
	}

	initScheduler(walkerInst, opts.ScrapeInterval)
	startServer(opts.Server)
}

func startServer(serverConf config.ServerConfiguration) {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>Elasticsearch Exporter</title></head>
			<body>
			<h1>Elasticsearch Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Errorf("failed handling writer: %s", err.Error())
		}
	})

	listen := serverConf.Listen + ":" + strconv.Itoa(serverConf.Port)
	var err error
	if serverConf.ServerCertificate != "" && serverConf.ServerKeyFile != "" {
		log.Infof("Starting HTTPS server on %s (Cert: %s, Key%s)", listen, serverConf.ServerCertificate, serverConf.ServerKeyFile)
		err = http.ListenAndServeTLS(listen, serverConf.ServerCertificate, serverConf.ServerKeyFile, r)
	} else {
		log.Infof("Starting HTTP server on %s", listen)
		err = http.ListenAndServe(listen, r)
	}

	if err != nil {
		log.Fatal("Could not start Server: ", err)
	}
}

func initScheduler(walkerInst walker.Walker, duration time.Duration) {

	log.Infof("Scheduling scrape every %s", duration.String())
	utils.Schedule(func() {
		err := walkerInst.Walk()
		if err != nil {
			log.Errorf("could not walk the specified path: ", err.Error())
		}
	}, duration)
}
