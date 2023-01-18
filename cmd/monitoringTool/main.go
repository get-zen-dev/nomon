package main

import (
	"log"

	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

func main() {

	stat := getServerStats.StartMonitoring()
	log.Printf("Cloudron stats:\n CPU Usage: %f\n Disk Usage: %f\n Memory Usage: %f\n", stat.CpuStatus, stat.DiskStatus, stat.MemStatus)
	// log.Println("Starting server on http://127.0.0.1:8080/")

	// r := newRouter()
	// log.Fatal(http.ListenAndServe(":8080", r))
}
