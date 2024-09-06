package parser

import (
	"strconv"
	"strings"

	"smartctl_ssacli_exporter/applog"	
	"github.com/go-kit/log/level"
)

func toINT(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		level.Error(applog.Logger).Log("Error parse int", err)
		return 0
	}
	return int64(i)
}

func toFLO(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		level.Error(applog.Logger).Log("Error parse float ", err)
		return 0.0
	}
	return float64(i)
}

func trim(s string) string {
	return strings.Trim(s, " \t")
}
