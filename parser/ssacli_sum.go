package parser

import (
	"strings"
)

// SsacliSum data structure for output
type SsacliSum struct {
	ContNumber    int
	SsacliSumData []SsacliSumData
}

// SsacliSumData data structure for output
type SsacliSumData struct {
	Slot           int64
	SlotID         string
	SerialNumber   string
	ContStatus     string
	FirmVersion    string
	TotalCacheSize float64
	AvailCacheSize float64
	BatteryStatus  string
	ContTemp       float64
	CahceModuTemp  float64
	BatteryTemp    float64
	Encryption     string
	DriverName     string
	DriverVersion  string
}

// ParseSsacliSum return specific metric
func ParseSsacliSum(s string) *SsacliSum {
	data := parseSmartAttrs(s)

	return data
}

func parseSmartAttrs(s string) *SsacliSum {

	var (
		contNumber int
		tmp        SsacliSumData
		// tmp1 ssacliSumData
	)

	for _, line := range strings.Split(s, "\n") {
		kvs := strings.Trim(line, " \t")
		kv := strings.Split(kvs, ": ")

		if len(kv) == 2 {

			switch kv[0] {
			case "Slot":
				tmp.Slot = toINT(kv[1])
				tmp.SlotID = kv[1]
			case "Serial Number":
				tmp.SerialNumber = kv[1]
			case "Controller Status":
				tmp.ContStatus = kv[1]
			case "Firmware Version":
				tmp.FirmVersion = kv[1]
			case "Total Cache Size":
				tmp.TotalCacheSize = toFLO(kv[1])
			case "Total Cache Memory Available":
				tmp.AvailCacheSize = toFLO(kv[1])
			case "Battery/Capacitor Status":
				tmp.BatteryStatus = kv[1]
			case "Controller Temperature (C)":
				tmp.ContTemp = toFLO(kv[1])
			case "Cache Module Temperature (C)":
				tmp.CahceModuTemp = toFLO(kv[1])
			case "Capacitor Temperature  (C)":
				tmp.BatteryTemp = toFLO(kv[1])
			case "Encryption":
				tmp.Encryption = kv[1]
			case "Driver Name":
				tmp.DriverName = kv[1]
			case "Driver Version":
				tmp.DriverVersion = kv[1]
			}
			// } else {
			// 	contNumber++
			// 	switch kv[0] {
			// 	case "Slot":
			// 		tmp1.slot = toINT(kv[1])
			// 	case "Serial Number":
			// 		tmp1.serialNumber = kv[1]
			// 	case "Controller Status":
			// 		tmp1.contStatus = kv[1]
			// 	case "Firmware Version":
			// 		tmp1.firmVersion = kv[1]
			// 	case "Total Cache Size":
			// 		tmp1.totalCacheSize = toFLO(kv[1])
			// 	case "Total Cache Memory Available":
			// 		tmp1.availCacheSize = toFLO(kv[1])
			// 	case "Battery/Capacitor Status":
			// 		tmp1.batteryStatus = kv[1]
			// 	case "Controller Temperature (C)":
			// 		tmp1.contTemp = toFLO(kv[1])
			// 	case "Cache Module Temperature (C)":
			// 		tmp1.cahceModuTemp = toFLO(kv[1])
			// 	case "Capacitor Temperature  (C)":
			// 		tmp1.batteryTemp = toFLO(kv[1])
			// 	case "Encryption":
			// 		tmp1.encryption = kv[1]
			// 	case "Driver Name":
			// 		tmp1.driverName = kv[1]
			// 	case "Driver Version":
			// 		tmp1.driverVersion = kv[1]
			// 	}

		}
	}

	data := SsacliSum{
		ContNumber: contNumber,
		SsacliSumData: []SsacliSumData{
			tmp,
			// tmp1,
		},
	}
	return &data
}
