package main

import (
	"os"
	"sync"

	"github.com/Setom29/CloudronMonitoring/pkg/cloudronStatus"
	"github.com/Setom29/CloudronMonitoring/pkg/parseData"
)

func main() {
	var wg sync.WaitGroup
	cfg, err := parseData.ParseConfig("./config/config.yml")
	if err != nil {
		os.Exit(1)
	}
	wg.Add(1)
	go cloudronStatus.StartMonitoring(cfg, &wg)
	wg.Wait()

}
