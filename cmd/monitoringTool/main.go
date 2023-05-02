package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	"github.com/Setom29/CloudronMonitoring/pkg/parseConfig"
	"github.com/Setom29/CloudronMonitoring/pkg/webInterface"
	log "github.com/sirupsen/logrus"
)

func setLogs(lvl int) {
	log.SetLevel(log.Level(lvl))
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func main() {
	log.Trace("main:main")
	args, report, err := parseConfig.Parse("./data/config.yml")
	if err != nil {
		log.Fatal("Error parsing config: ", err)
		return
	}
	setLogs(args.MonitorLogLevel)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	m := monitor.NewMonitor(args, report)
	log.Debug(m.Args)
	go m.StartMonitoring(sigChan)
	http.HandleFunc("/", webInterface.MakeIndexHandler(m))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", m.Args.Port), nil))
}
