package flags

import (
	"errors"
	"flag"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
)

func ParseFlags() (monitor.Args, error) {
	f := monitor.Args{}

	limit := flag.Float64("p", 90, "max usage in % of RAM, CPU and Disk")
	duration := flag.Int("t", 300, "time period for checking")
	checkTime := flag.Int("c", 30, "check server status every X seconds")
	dbFile := flag.String("f", "serverStats.db", "database filename")
	port := flag.String("p", "8080", "port")

	flag.Parse()

	if *limit > 100 || *limit < 0 {
		return monitor.Args{}, errors.New("wrong value for persentage")
	}
	if *duration < 0 {
		return monitor.Args{}, errors.New("wrong value for duration")
	}
	if *checkTime < 0 || *checkTime > *duration {
		return monitor.Args{}, errors.New("wrong value for check time")
	}

	f.Limit = *limit
	f.Duration = *duration
	f.CheckTime = *checkTime
	f.DBFile = *dbFile
	f.Port = *port
	return f, nil
}
