package parser

import (
	"regexp"
	"strings"
)

// SmartctlDisk data structure for output
type SmartctlDisk struct {
	SmartctlDiskDataInfo []SmartctlDiskDataInfo
	SmartctlDiskDataAttr []SmartctlDiskDataAttr
}

// SmartctlDiskDataInfo comment
type SmartctlDiskDataInfo struct {
	Model    string
	SN       string
	RotRate  string
	FromFact string
}

// SmartctlDiskDataAttr comment
type SmartctlDiskDataAttr struct {
	RawReadErrorRate      float64
	ReallocatedSectorCt   float64
	PowerOnHours          float64
	PowerCycleCount       float64
	RuntimeBadBlock       float64
	EndToEndError         float64
	ReportedUncorrect     float64
	CommandTimeout        float64
	HardwareECCRecovered  float64
	ReallocatedEventCount float64
	CurrentPendingSector  float64
	OfflineUncorrectable  float64
	UDMACRCErrorCount     float64
	UnusedRsvdBlkCntTot   float64
}

// ParseSmartctlDisk return specific metric
func ParseSmartctlDisk(s string) *SmartctlDisk {

	dataAtr := SmartctlDiskDataAttr{}
	dataInfo := SmartctlDiskDataInfo{}
	for _, section := range strings.Split(s, "=== START OF ") {
		if strings.Index(section, "INFORMATION SECTION ===") > -1 {
			dataInfo = parseSmartctlDiskInfo(section)
		} else if strings.Index(section, "READ SMART DATA SECTION ===") > -1 {
			dataAtr = parseSmartctlDiskAtr(section)
		}
	}

	data := SmartctlDisk{
		SmartctlDiskDataAttr: []SmartctlDiskDataAttr{
			dataAtr,
		},
		SmartctlDiskDataInfo: []SmartctlDiskDataInfo{
			dataInfo,
		},
	}

	return &data
}

func parseSmartctlDiskInfo(s string) SmartctlDiskDataInfo {

	var (
		tmp SmartctlDiskDataInfo
	)

	for _, line := range strings.Split(s, "\n") {
		kvs := strings.Trim(line, " \t")
		kv := strings.Split(kvs, ": ")

		if len(kv) == 2 {
			switch kv[0] {
			case "Device Model":
				tmp.Model = trim(kv[1])
			case "Serial Number":
				tmp.SN = trim(kv[1])
			case "Rotation Rate":
				tmp.RotRate = trim(kv[1])
			case "Form Factor":
				tmp.FromFact = trim(kv[1])
			}
		}
	}

	return tmp
}

func parseSmartctlDiskAtr(s string) SmartctlDiskDataAttr {

	var (
		tmp SmartctlDiskDataAttr
	)

	reSpaces := regexp.MustCompile(`\s+`)
	lines := strings.Split(s, "\n")

	for _, line := range lines {
		vals := reSpaces.Split(trim(line), -1)

		if len(vals) < 10 {
			continue
		}

		switch vals[1] {
		case "Raw_Read_Error_Rate":
			tmp.RawReadErrorRate = toFLO(vals[9])
		case "Reallocated_Sector_Ct":
			tmp.ReallocatedSectorCt = toFLO(vals[9])
		case "Power_On_Hours":
			tmp.PowerOnHours = toFLO(vals[9])
		case "Power_Cycle_Count":
			tmp.PowerCycleCount = toFLO(vals[9])
		case "Runtime_Bad_Block":
			tmp.RuntimeBadBlock = toFLO(vals[9])
		case "End-to-End_Error":
			tmp.EndToEndError = toFLO(vals[9])
		case "Reported_Uncorrect":
			tmp.ReportedUncorrect = toFLO(vals[9])
		case "Command_Timeout":
			tmp.CommandTimeout = toFLO(vals[9])
		case "Hardware_ECC_Recovered":
			tmp.HardwareECCRecovered = toFLO(vals[9])
		case "Reallocated_Event_Count":
			tmp.ReallocatedEventCount = toFLO(vals[9])
		case "Current_Pending_Sector":
			tmp.CurrentPendingSector = toFLO(vals[9])
		case "Offline_Uncorrectable":
			tmp.OfflineUncorrectable = toFLO(vals[9])
		case "UDMA_CRC_Error_Count":
			tmp.UDMACRCErrorCount = toFLO(vals[9])
		case "Unused_Rsvd_Blk_Cnt_Tot":
			tmp.UnusedRsvdBlkCntTot = toFLO(vals[9])
		}
	}

	return tmp
}
