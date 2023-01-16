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
	// serverStatus.getDisk()

	log.Println("CPU status:", serverStatus.CpuStatus)
	log.Println("Memory status:", serverStatus.MemStatus)
}

func (serverStatus *CloudronServerStatus) getCpu() {
	before, err := cpu.Get()
	if err != nil {
		log.Println("Error getting CPU: ", err)
	}
	time.Sleep(time.Duration(serverStatus.duration) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		log.Println("Error getting CPU: ", err)
	}
	total := float64(after.Total - before.Total)
	serverStatus.MemStatus = total
}

func (serverStatus *CloudronServerStatus) getMem() {
	before, err := memory.Get()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}
	time.Sleep(time.Duration(serverStatus.duration) * time.Second)
	after, err := memory.Get()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}
	total := float64(after.Total - before.Total)
	serverStatus.CpuStatus = total
}

// func (serverStatus *CloudronServerStatus) getDisk() (float64, error) {
// 	before, err := disk.Get()
// 	if err != nil {
// 		return float64(0), err
// 	}
// 	time.Sleep(time.Duration(serverStatus.duration) * time.Second)
// 	after, err := disk.Get()
// 	if err != nil {
// 		return float64(0), err
// 	}
// 	total := float64(after.Total - before.Total)
// 	return float64(after.Used-before.Used) / total * 100, nil
// }

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
