package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

type CloudronServerStatus struct {
	CpuStatus   float64
	MemStatus   float64
	DiskStatus  float64
	LastUpdated time.Time
	duration    time.Duration
	err         error
}

func NewCloudronServerConnection(duration uint64) *CloudronServerStatus {
	return &CloudronServerStatus{CpuStatus: 0, MemStatus: 0, DiskStatus: 0, LastUpdated: time.Now(), err: nil, duration: time.Duration(duration)}
}

func main() {
	StartMonitoring()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func StartMonitoring() {
	serverStatus := NewCloudronServerConnection(5)
	serverStatus.getCpu()
	serverStatus.getMem()
	serverStatus.getDisk()

	log.Println("CPU status:", serverStatus.CpuStatus)
	log.Println("Memory status:", serverStatus.MemStatus)
}

func (serverStatus *CloudronServerStatus) getCpu() {
	cpuInfo, err := cpu.Get()
	if err != nil {
		log.Println("Error getting CPU: ", err)
	}
	serverStatus.MemStatus = float64(cpuInfo.System) / float64(cpuInfo.Total) * 100
}

func (serverStatus *CloudronServerStatus) getMem() {
	memInfo, err := memory.Get()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}

	serverStatus.CpuStatus = float64(memInfo.Used) / float64(memInfo.Total) * 100
}

func (serverStatus *CloudronServerStatus) getDisk() {
	//not implemented
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
