package flags

import (
	"errors"
	"flag"
	"time"
)

type Flags struct {
	Limit     float64
	Duration  time.Duration
	CheckTime time.Duration
	DBFile    string
}

func ParseFlags() (Flags, error) {
	f := Flags{}

	limit := flag.Float64("p", 90, "max usage in %% of RAM, CPU and Disk")
	duration := flag.Duration("t", time.Second*60, "time period for checking")
	checkTime := flag.Duration("c", time.Second*10, "check server status every _ seconds")
	dbFile := flag.String("f", "serverStats.db", "database filename")

	flag.Parse()

	if *limit > 100 || *limit < 0 {
		return Flags{}, errors.New("wrong value for persentage")
	}
	if *duration < time.Second*0 {
		return Flags{}, errors.New("wrong value for duration")
	}
	if *checkTime < time.Second*0 || *checkTime > *duration {
		return Flags{}, errors.New("wrong value for check time")
	}

	f.Limit = *limit
	f.Duration = *duration
	f.CheckTime = *checkTime
	f.DBFile = *dbFile
	return f, nil
}
