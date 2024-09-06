package collector

import (
	"os/exec"

	"smartctl_ssacli_exporter/parser"
	"smartctl_ssacli_exporter/applog"
	
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// ConID save controller slot number
var ConID string

var _ prometheus.Collector = &SsacliSumCollector{}

// SsacliSumCollector Contain raid controller detail information
type SsacliSumCollector struct {
	id                 string
	hwConSlotDesc      *prometheus.Desc
	cacheSizeDesc      *prometheus.Desc
	availCacheSizeDesc *prometheus.Desc
	hwConTempDesc      *prometheus.Desc
	cahceModuTempDesc  *prometheus.Desc
	batteryTempDesc    *prometheus.Desc
}

// NewSsacliSumCollector Create new collector
func NewSsacliSumCollector() *SsacliSumCollector {
	// Init labels
	var (
		namespace = "ssacli"
		subsystem = "hw_raid_controller"
		labels    = []string{
			"raidControllerSN",
			"raidControllerStatus",
			"raidControllerFirmVersion",
			"raidControllerBatteryStatus",
			"raidControllerEncryption",
			"raidControllerDriverName",
			"raidControllerDriverVersion",
		}
	)
	// Rerutn Colected metric to ch <-
	// Include labels
	return &SsacliSumCollector{
		hwConSlotDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "slot"),
			"Hardware raid controller slot usage",
			labels,
			nil,
		),
		cacheSizeDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cacheSize"),
			"Hardware raid controller total cahce size",
			labels,
			nil,
		),
		availCacheSizeDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "available_cacheSize"),
			"Hardware raid controller total available cahce size",
			labels,
			nil,
		),
		hwConTempDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "temperature"),
			"Hardware raid controller hardware temperature",
			labels,
			nil,
		),
		cahceModuTempDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "temperature_cacheModule"),
			"Hardware raid controller cache module temperature",
			labels,
			nil,
		),
		batteryTempDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "temperature_battery"),
			"Hardware raid controller battery/capacitor module temperature",
			labels,
			nil,
		),
	}
}

// Describe return all description to chanel
func (c *SsacliSumCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.hwConSlotDesc,
		c.cacheSizeDesc,
		c.availCacheSizeDesc,
		c.hwConTempDesc,
		c.cahceModuTempDesc,
		c.batteryTempDesc,
	}
	for _, d := range ds {
		ch <- d
	}
}

// Collect create collector
// Get metric
// Handle error
func (c *SsacliSumCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		level.Debug(applog.Logger).Log("Failed collecting metric", desc)
		level.Error(applog.Logger).Log("Failed collecting metric", err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}

func (c *SsacliSumCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	cmd := "ssacli ctrl all show detail"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		level.Debug(applog.Logger).Log("Smart log:", out)
		level.Error(applog.Logger).Log("smart log:", err)
		return nil, err
	}

	data := parser.ParseSsacliSum(string(out))

	if data == nil {
		level.Debug(applog.Logger).Log("ssacli parse sum:", data)
		level.Error(applog.Logger).Log("Unable parse data from ssacli sumarry exporter", "No parsed data")
		return nil, nil
	}

	for i := range data.SsacliSumData {

		var (
			labels = []string{
				data.SsacliSumData[i].SerialNumber,
				data.SsacliSumData[i].ContStatus,
				data.SsacliSumData[i].FirmVersion,
				data.SsacliSumData[i].BatteryStatus,
				data.SsacliSumData[i].Encryption,
				data.SsacliSumData[i].DriverName,
				data.SsacliSumData[i].DriverVersion,
			}
		)

		ConID = data.SsacliSumData[i].SlotID

		ch <- prometheus.MustNewConstMetric(
			c.hwConSlotDesc,
			prometheus.GaugeValue,
			float64(data.SsacliSumData[i].Slot),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.cacheSizeDesc,
			prometheus.GaugeValue,
			float64(data.SsacliSumData[i].TotalCacheSize),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.availCacheSizeDesc,
			prometheus.GaugeValue,
			float64(data.SsacliSumData[i].AvailCacheSize),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.hwConTempDesc,
			prometheus.GaugeValue,
			float64(data.SsacliSumData[i].ContTemp),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.cahceModuTempDesc,
			prometheus.GaugeValue,
			float64(data.SsacliSumData[i].CahceModuTemp),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.batteryTempDesc,
			prometheus.GaugeValue,
			float64(data.SsacliSumData[i].BatteryTemp),
			labels...,
		)

	}
	return nil, nil
}
