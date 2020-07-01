package parser

import (
	"strings"
)

// SsacliPhysDisk data structure for output
type SsacliPhysDisk struct {
	SsacliPhysDiskData []SsacliPhysDiskData
}

// SsacliPhysDiskData data structure for output
type SsacliPhysDiskData struct {
	Bay       string
	Status    string
	DriveType string
	IntType   string
	Size      string
	BlockSize string
	SN        string
	WWID      string
	CurTemp   float64
	MaxTemp   float64
	Model     string
}

// ParseSsacliPhysDisk return specific metric
func ParseSsacliPhysDisk(s string) *SsacliPhysDisk {
	data := parseSsacliPhysDisk(s)

	return data
}

func parseSsacliPhysDisk(s string) *SsacliPhysDisk {

	var (
		tmp SsacliPhysDiskData
	)
	for _, line := range strings.Split(s, "\n") {
		kvs := strings.Trim(line, " \t")
		kv := strings.Split(kvs, ": ")

		if len(kv) == 2 {

			switch kv[0] {
			case "Bay":
				tmp.Bay = kv[1]
			case "Serial Number":
				tmp.SN = kv[1]
			case "Status":
				tmp.Status = kv[1]
			case "Drive Type":
				tmp.DriveType = kv[1]
			case "Interface Type":
				tmp.IntType = kv[1]
			case "Size":
				tmp.Size = kv[1]
			case "Logical/Physical Block Size":
				tmp.BlockSize = kv[1]
			case "WWID":
				tmp.WWID = kv[1]
			case "Model":
				tmp.Model = kv[1]
			case "Current Temperature (C)":
				tmp.CurTemp = toFLO(kv[1])
			case "Maximum Temperature (C)":
				tmp.MaxTemp = toFLO(kv[1])
			}
		}
	}

	data := SsacliPhysDisk{
		SsacliPhysDiskData: []SsacliPhysDiskData{
			tmp,
		},
	}
	return &data
}
