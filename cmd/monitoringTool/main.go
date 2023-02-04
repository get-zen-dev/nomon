package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Flags struct {
	limit    int
	duration time.Duration
	dbFile   string
}

func parseFlags() (flags, error) {
	f := flags{}
	limit := flag.Int("p", 90, "max usage in %% of RAM, CPU and Disk")
	duration := flag.Duration("t", time.Duration(60), "period of time in seconds")
	dbFile := flag.String("f", "serverStats.db", "database filename")
	flag.Parse()
	if *limit > 100 || *limit < 0 {
		return flags{}, errors.New("wrong value for persentage")
	}
	if *duration < time.Duration(0) {

		return flags{}, errors.New("wrong value for duration")
	}
	f.limit = *limit
	f.duration = *duration
	f.dbFile = *dbFile
	return f, nil

}

func main() {
	dbFile := "file"

	f, err := parseFlags()
	if err != nil {
		log.Println(err)
		return
	}

	// log.Println("Starting server on http://127.0.0.1:8080/")

	// r := newRouter()
	// log.Fatal(http.ListenAndServe(":8080", r))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go monitor.StartMonitor(sigChan, dbFile)
}
