package main

import "github.com/Setom29/CloudronMonitoring/pkg/cloudronStatus"

func main() {
	cloudronStatus.StartMonitoring("./config/config.yml")

}
