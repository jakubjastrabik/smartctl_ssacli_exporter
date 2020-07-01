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
	collector.NewSsacliPhysDiskCollector("").Describe(ch)
	collector.NewSmartctlDiskCollector("", 0).Describe(ch)
}

// Collect sends the collected metrics from each of the collectors to
// exporter.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	collector.NewSsacliSumCollector().Collect(ch)

	// Inicializacia novych collectorov pre kazdy prikaz samostatne.
	// V provom krouku zistime kolko mame vlastne diskov.
	// Nasledne kazdy disk dostane svojho workera ktory zisti
	// zvysne info tym ze zavola specificky collector.

	cmd := "ssacli ctrl slot=0 pd all show status| grep . | cut -f5 -d' '"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		log.Printf("[ERROR] failed collecting metric %v: %v", out, err)
		return
	}

	physDisk := strings.Split(string(out), "\n")
	physDiskN := 1
	for _, physDisk := range physDisk {
		if physDisk != "" {
			collector.NewSsacliPhysDiskCollector(physDisk).Collect(ch)
			collector.NewSmartctlDiskCollector(physDisk, physDiskN).Collect(ch)
			physDiskN++
		}
	}

	// Export logic raid status

	cmd = "ssacli ctrl slot=0 ld all show status| grep . | cut -f5 -d' '"
	out, err = exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		log.Printf("[ERROR] failed collecting metric %v: %v", out, err)
		return
	}

	logDisk := strings.Split(string(out), "\n")
	for _, logDisk := range logDisk {
		if logDisk != "" {
			collector.NewSsacliLogDiskCollector(logDisk).Collect(ch)
		}
	}
}
