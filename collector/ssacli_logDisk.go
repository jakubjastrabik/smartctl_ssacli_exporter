package collector

import (
	"os/exec"

	"smartctl_ssacli_exporter/parser"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &SsacliLogDiskCollector{}

// SsacliLogDiskCollector Contain raid controller detail information
type SsacliLogDiskCollector struct {
	diskID    string
	conID     string
	cylinders *prometheus.Desc
}

// NewSsacliLogDiskCollector Create new collector
func NewSsacliLogDiskCollector(diskID, conID string) *SsacliLogDiskCollector {
	// Init labels
	var (
		namespace = "ssacli"
		subsystem = "logical_array"
		labels    = []string{
			"Size",
			"Status",
			"Caching",
			"UID",
			"LName",
			"LID",
		}
	)

	// Rerutn Colected metric to ch <-
	// Include labels
	return &SsacliLogDiskCollector{
		diskID: diskID,
		conID:  conID,
		cylinders: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cylinders"),
			"Logical array cylinder count",
			labels,
			nil,
		),
	}
}

// Describe return all description to chanel
func (c *SsacliLogDiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.cylinders,
	}
	for _, d := range ds {
		ch <- d
	}
}

// Collect create collector
// Get metric
// Handle error
func (c *SsacliLogDiskCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		// log.Debugln("[ERROR] failed collecting metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}

func (c *SsacliLogDiskCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	if c.diskID == "" {
		return nil, nil
	}

	cmd := "ssacli ctrl slot=" + c.conID + " ld " + c.diskID + " show | grep ."
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		// log.Debugln("[ERROR] smart log: \n%s\n", out)
		return nil, err
	}

	data := parser.ParseSsacliLogDisk(string(out))

	if data == nil {
		// log.Fatal("Unable get data from ssacli logical array exporter")
		return nil, nil
	}

	var (
		labels = []string{
			data.SsacliLogDiskData[0].Size,
			data.SsacliLogDiskData[0].Status,
			data.SsacliLogDiskData[0].Caching,
			data.SsacliLogDiskData[0].UID,
			data.SsacliLogDiskData[0].LName,
			data.SsacliLogDiskData[0].LID,
		}
	)

	ch <- prometheus.MustNewConstMetric(
		c.cylinders,
		prometheus.GaugeValue,
		float64(data.SsacliLogDiskData[0].Cylinders),
		labels...,
	)

	return nil, nil
}
