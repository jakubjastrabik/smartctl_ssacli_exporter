package main

import (
	"flag"
	"net/http"

	"smartctl_ssacli_exporter/exporter"
	"smartctl_ssacli_exporter/applog"
	
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddr  = flag.String("listen", ":9633", "address for exporter")
	metricsPath = flag.String("path", "/metrics", "URL path for surfacing collected metrics")
	logLevel = flag.String("log-level", "debug", "Logging level (debug, info, warn or error")
)

func main() {
	// Parse command-line flags and set up logging level
	flag.Parse()

	// Initialize logger with specified log level
	applog.ApplogInit(*logLevel)

	// Register metrics with Prometheus	
	prometheus.MustRegister(exporter.New())

	// Serve the metrics on the specified path.
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	level.Info(applog.Logger).Log("msg", "Beginning to serve on port " + *listenAddr)

	// Start the server and handle any errors that may occur.
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		level.Error(applog.Logger).Log("Cannot start exporter", err)
		panic(err) // Crash and burn if we cannot start the server.
	}
}
