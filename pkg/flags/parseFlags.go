package flags

import (
	"errors"
	"flag"
)

type Flags struct {
	Limit     float64
	Duration  int
	CheckTime int
	DBFile    string
}

func ParseFlags() (Flags, error) {
	f := Flags{}

	limit := flag.Float64("p", 90, "max usage in %% of RAM, CPU and Disk")
	duration := flag.Int("t", 300, "time period for checking")
	checkTime := flag.Int("c", 30, "check server status every _ seconds")
	dbFile := flag.String("f", "serverStats.db", "database filename")

	flag.Parse()

	if *limit > 100 || *limit < 0 {
		return Flags{}, errors.New("wrong value for persentage")
	}
	if *duration < 0 {
		return Flags{}, errors.New("wrong value for duration")
	}
	if *checkTime < 0 || *checkTime > *duration {
		return Flags{}, errors.New("wrong value for check time")
	}

	f.Limit = *limit
	f.Duration = *duration
	f.CheckTime = *checkTime
	f.DBFile = *dbFile
	return f, nil
}
