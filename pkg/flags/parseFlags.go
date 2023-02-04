package flags

import (
	"errors"
	"flag"
	"time"
)

type Flags struct {
	limit    float64
	duration time.Duration
	dbFile   string
}

func parseFlags() (Flags, error) {
	f := Flags{}

	limit := flag.Float64("p", 90, "max usage in %% of RAM, CPU and Disk")
	duration := flag.Duration("t", time.Duration(60), "period of time in seconds (min 60 seconds)")
	dbFile := flag.String("f", "serverStats.db", "database filename")

	flag.Parse()
	if *limit > 100 || *limit < 0 {
		return Flags{}, errors.New("wrong value for persentage")
	}
	if *duration < time.Duration(0) {

		return Flags{}, errors.New("wrong value for duration")
	}
	f.limit = *limit
	f.duration = *duration
	f.dbFile = *dbFile
	return f, nil

}
