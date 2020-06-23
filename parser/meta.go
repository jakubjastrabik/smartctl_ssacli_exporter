package parser

import (
	"log"
	"strconv"
)

func toINT(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalln(err)
	}
	return int64(i)
}

func toFLO(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return float64(i)
}
