package parser

import (
	"strings"
)

// SsacliLogDisk data structure for output
type SsacliLogDisk struct {
	SsacliLogDiskData []SsacliLogDiskData
}

// SsacliLogDiskData data structure for output
type SsacliLogDiskData struct {
	Size      string
	Cylinders float64
	Status    string
	Caching   string
	UID       string
	LName     string
	LID       string
}

// ParseSsacliLogDisk return specific metric
func ParseSsacliLogDisk(s string) *SsacliLogDisk {
	data := parseSsacliLogDisk(s)

	return data
}

func parseSsacliLogDisk(s string) *SsacliLogDisk {

	var (
		tmp SsacliLogDiskData
	)

	for _, line := range strings.Split(s, "\n") {
		kvs := strings.Trim(line, " \t")
		kv := strings.Split(kvs, ": ")

		if len(kv) == 2 {

			switch kv[0] {
			case "Size":
				tmp.Size = kv[1]
			case "Cylinders":
				tmp.Cylinders = toFLO(kv[1])
			case "Status":
				tmp.Status = kv[1]
			case "Caching":
				tmp.Caching = kv[1]
			case "Unique Identifier":
				tmp.UID = kv[1]
			case "Disk Name":
				tmp.LName = kv[1]
			case "Logical Drive Label":
				tmp.LID = kv[1]
			}
		}
	}

	data := SsacliLogDisk{
		SsacliLogDiskData: []SsacliLogDiskData{
			tmp,
		},
	}
	return &data
}
