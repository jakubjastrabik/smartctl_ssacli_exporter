package collector

import (
	"os/exec"
	"strconv"

	"smartctl_ssacli_exporter/parser"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &SmartctlDiskCollector{}

// SmartctlDiskCollector Contain raid controller detail information
type SmartctlDiskCollector struct {
	diskID string
	diskN  int

	model    *prometheus.Desc
	sn       *prometheus.Desc
	rotRate  *prometheus.Desc
	fromFact *prometheus.Desc

	rawReadErrorRate      *prometheus.Desc
	reallocatedSectorCt   *prometheus.Desc
	powerOnHours          *prometheus.Desc
	powerCycleCount       *prometheus.Desc
	runtimeBadBlock       *prometheus.Desc
	endToEndError         *prometheus.Desc
	reportedUncorrect     *prometheus.Desc
	commandTimeout        *prometheus.Desc
	hardwareECCRecovered  *prometheus.Desc
	reallocatedEventCount *prometheus.Desc
	currentPendingSector  *prometheus.Desc
	offlineUncorrectable  *prometheus.Desc
	uDMACRCErrorCount     *prometheus.Desc
	unusedRsvdBlkCntTot   *prometheus.Desc
}

// NewSmartctlDiskCollector Create new collector
func NewSmartctlDiskCollector(diskID string, diskN int) *SmartctlDiskCollector {
	// Init labels
	var (
		namespace = "smartctl"
		subsystem = "physical_disk"
		labels    = []string{
			"diskID",
			"model",
			"sn",
			"rotRate",
			"fromFact",
		}
	)

	// Rerutn Colected metric to ch <-
	// Include labels
	return &SmartctlDiskCollector{
		diskID: diskID,
		diskN:  diskN,
		rawReadErrorRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "rawReadErrorRate"),
			"Smartctl raw read error rate",
			labels,
			nil,
		),
		reallocatedSectorCt: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reallocatedSectorCt"),
			"Smartctl reallocated sector ct",
			labels,
			nil,
		),
		powerOnHours: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "powerOnHours"),
			"Smartctl power on hours",
			labels,
			nil,
		),
		powerCycleCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "powerCycleCount"),
			"Smartctl power cycle down count",
			labels,
			nil,
		),
		runtimeBadBlock: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "runtimeBadBlock"),
			"Smartctl runtime bad block",
			labels,
			nil,
		),
		endToEndError: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "endToEndError"),
			"Smartctl end to end error",
			labels,
			nil,
		),
		reportedUncorrect: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reportedUncorrect"),
			"Smartctl reported uncorrect",
			labels,
			nil,
		),
		commandTimeout: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "commandTimeout"),
			"Smartctl command timeout",
			labels,
			nil,
		),
		hardwareECCRecovered: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "hardwareECCRecovered"),
			"Smartctl hardware ecc recovered",
			labels,
			nil,
		),
		reallocatedEventCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reallocatedEventCount"),
			"Smartctl reallocated event count",
			labels,
			nil,
		),
		currentPendingSector: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "currentPendingSector"),
			"Smartctl current pending sector",
			labels,
			nil,
		),
		offlineUncorrectable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "offlineUncorrectable"),
			"Smartctl offline uncorrectable",
			labels,
			nil,
		),
		uDMACRCErrorCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "uDMACRCErrorCount"),
			"Smartctl ud macrc error count",
			labels,
			nil,
		),
		unusedRsvdBlkCntTot: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "unusedRsvdBlkCntTot"),
			"Smartctl unused rsvd block Count Total",
			labels,
			nil,
		),
	}
}

// Describe return all description to chanel
func (c *SmartctlDiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.rawReadErrorRate,
		c.reallocatedSectorCt,
		c.powerOnHours,
		c.powerCycleCount,
		c.runtimeBadBlock,
		c.endToEndError,
		c.reportedUncorrect,
		c.commandTimeout,
		c.hardwareECCRecovered,
		c.reallocatedEventCount,
		c.currentPendingSector,
		c.offlineUncorrectable,
		c.uDMACRCErrorCount,
		c.unusedRsvdBlkCntTot,
	}
	for _, d := range ds {
		ch <- d
	}
}

// Collect create collector
// Get metric
// Handle error
func (c *SmartctlDiskCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		// log.Debugln("[ERROR] failed collecting metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}

func (c *SmartctlDiskCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	if c.diskID == "" {
		return nil, nil
	}

	cmd := "smartctl -iA -d cciss," + strconv.Itoa(c.diskN) + " /dev/sda | grep ."
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		// log.Debugln("[ERROR] smart log: \n%s\n", out)
		return nil, err
	}
	data := parser.ParseSmartctlDisk(string(out))

	if data == nil {
		// log.Fatal("Unable get data from smartctl exporter")
		return nil, nil
	}

	var (
		labels = []string{
			c.diskID,
			data.SmartctlDiskDataInfo[0].Model,
			data.SmartctlDiskDataInfo[0].SN,
			data.SmartctlDiskDataInfo[0].RotRate,
			data.SmartctlDiskDataInfo[0].FromFact,
		}
	)

	ch <- prometheus.MustNewConstMetric(
		c.rawReadErrorRate,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].RawReadErrorRate),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.reallocatedSectorCt,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].ReallocatedSectorCt),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.powerOnHours,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].PowerOnHours),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.powerCycleCount,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].PowerCycleCount),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.runtimeBadBlock,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].RuntimeBadBlock),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.endToEndError,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].EndToEndError),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.reportedUncorrect,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].ReportedUncorrect),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.commandTimeout,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].CommandTimeout),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.hardwareECCRecovered,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].HardwareECCRecovered),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.reallocatedEventCount,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].ReallocatedEventCount),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.currentPendingSector,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].CurrentPendingSector),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.offlineUncorrectable,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].OfflineUncorrectable),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.uDMACRCErrorCount,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].UDMACRCErrorCount),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.unusedRsvdBlkCntTot,
		prometheus.GaugeValue,
		float64(data.SmartctlDiskDataAttr[0].UnusedRsvdBlkCntTot),
		labels...,
	)

	return nil, nil
}
