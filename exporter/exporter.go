package exporter

import (
	"log"
	"os/exec"
	"strings"

	"github.com/jakubjastrabik/smartctl_ssacli_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

// An Exporter is a Prometheus exporter for metrics.
// It wraps all metrics collectors and provides a single global
// exporter which can serve metrics.
//
// It implements the exporter.Collector interface in order to register
// with Prometheus.
type Exporter struct {
}

var _ prometheus.Collector = &Exporter{}

// New creates a new Exporter which collects metrics by creating a apcupsd
// client using the input ClientFunc.
func New() *Exporter {
	return &Exporter{}
}

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	collector.NewSsacliSumCollector().Describe(ch)
	collector.NewSsacliPhysDiskCollector("", "").Describe(ch)
	collector.NewSmartctlDiskCollector("", 0).Describe(ch)
	collector.NewSsacliLogDiskCollector("", "").Describe(ch)
}

// Collect sends the collected metrics from each of the collectors to
// exporter.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	collector.NewSsacliSumCollector().Collect(ch)
	conID := collector.ConID

	cmd := "ssacli ctrl slot=" + conID + " pd all show status| grep . | cut -f5 -d' '"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		log.Printf("[ERROR] failed collecting metric %v: %v", out, err)
		return
	}

	physDisk := strings.Split(string(out), "\n")
	physDiskN := 1
	for _, physDisk := range physDisk {
		if physDisk != "" {
			collector.NewSsacliPhysDiskCollector(physDisk, conID).Collect(ch)
			collector.NewSmartctlDiskCollector(physDisk, physDiskN).Collect(ch)
			physDiskN++
		}
	}

	// Export logic raid status

	cmd = "ssacli ctrl slot=" + conID + " ld all show status| grep . | cut -f5 -d' '"
	out, err = exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		log.Printf("[ERROR] failed collecting metric %v: %v", out, err)
		return
	}

	logDisk := strings.Split(string(out), "\n")
	for _, logDisk := range logDisk {
		if logDisk != "" {
			collector.NewSsacliLogDiskCollector(logDisk, conID).Collect(ch)
		}
	}
}
