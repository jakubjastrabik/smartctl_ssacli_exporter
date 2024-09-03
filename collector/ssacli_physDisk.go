package collector

import (
	"os/exec"

	"smartctl_ssacli_exporter/parser"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &SsacliPhysDiskCollector{}

// SsacliPhysDiskCollector Contain raid controller detail information
type SsacliPhysDiskCollector struct {
	diskID  string
	conID   string
	curTemp *prometheus.Desc
	maxTemp *prometheus.Desc
}

// NewSsacliPhysDiskCollector Create new collector
func NewSsacliPhysDiskCollector(diskID, conID string) *SsacliPhysDiskCollector {
	// Init labels
	var (
		namespace = "ssacli"
		subsystem = "physical_disk"
		labels    = []string{
			"diskID",
			"Status",
			"DriveType",
			"IntType",
			"Size",
			"BlockSize",
			"SN",
			"WWID",
			"Model",
			"Bay",
		}
	)

	// Rerutn Colected metric to ch <-
	// Include labels
	return &SsacliPhysDiskCollector{
		diskID: diskID,
		conID:  conID,
		curTemp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "curTemp"),
			"Actual physical disk temperature",
			labels,
			nil,
		),
		maxTemp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "maxTmp"),
			"Physical disk maximum temperature",
			labels,
			nil,
		),
	}
}

// Describe return all description to chanel
func (c *SsacliPhysDiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.curTemp,
		c.maxTemp,
	}
	for _, d := range ds {
		ch <- d
	}
}

// Collect create collector
// Get metric
// Handle error
func (c *SsacliPhysDiskCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		// log.Debugln("[ERROR] failed collecting metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}

func (c *SsacliPhysDiskCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	if c.diskID == "" {
		return nil, nil
	}

	cmd := "ssacli ctrl slot=" + c.conID + " pd " + c.diskID + " show detail | grep ."
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		// log.Debugln("[ERROR] smart log: \n%s\n", out)
		return nil, err
	}

	data := parser.ParseSsacliPhysDisk(string(out))

	if data == nil {
		// log.Fatal("Unable get data from ssacli sumarry exporter")
		return nil, nil
	}

	var (
		labels = []string{
			c.diskID,
			data.SsacliPhysDiskData[0].Status,
			data.SsacliPhysDiskData[0].DriveType,
			data.SsacliPhysDiskData[0].IntType,
			data.SsacliPhysDiskData[0].Size,
			data.SsacliPhysDiskData[0].BlockSize,
			data.SsacliPhysDiskData[0].SN,
			data.SsacliPhysDiskData[0].WWID,
			data.SsacliPhysDiskData[0].Model,
			data.SsacliPhysDiskData[0].Bay,
		}
	)

	ch <- prometheus.MustNewConstMetric(
		c.curTemp,
		prometheus.GaugeValue,
		float64(data.SsacliPhysDiskData[0].CurTemp),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.maxTemp,
		prometheus.GaugeValue,
		float64(data.SsacliPhysDiskData[0].MaxTemp),
		labels...,
	)
	return nil, nil
}
