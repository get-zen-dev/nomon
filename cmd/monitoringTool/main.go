package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

type flags struct {
	limit    int
	duration time.Duration
}

func parseFlags() (flags, error) {
	f := flags{}
	limit := flag.Int("p", 90, "max percentage for RAM, CPU and Disk")
	duration := flag.Duration("t", time.Duration(300), "duration in seconds")
	flag.Parse()
	if *limit > 100 || *limit < 0 {
		return flags{}, errors.New("wrong value for persentage")
	}
	if *duration < time.Duration(0) {

		return flags{}, errors.New("wrong value for duration")
	}
	f.limit = *limit
	f.duration = *duration
	return f, nil

}

func main() {
	f, err := parseFlags()
	if err != nil {
		log.Println(err)
		return
	}
	getServerStats.StartMonitoring(f.limit, f.duration)
	// log.Println("Starting server on http://127.0.0.1:8080/")

	// r := newRouter()
	// log.Fatal(http.ListenAndServe(":8080", r))
}
