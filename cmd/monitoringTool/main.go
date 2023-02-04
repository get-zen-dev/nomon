package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
)

func main() {
	f, err := flags.parseFlags()
	if err != nil {
		log.Println(err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	m := monitor.NewMonitor(f)
	m.StartMonitoring(sigChan)

	// log.Println("Starting server on http://127.0.0.1:8080/")

	// r := newRouter()
	// log.Fatal(http.ListenAndServe(":8080", r))
}
