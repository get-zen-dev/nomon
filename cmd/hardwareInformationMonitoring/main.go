package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type CloudronServerStatus struct {
	CpuStatus   float64
	MemStatus   float64
	DiskStatus  float64
	LastUpdated time.Time
	duration    time.Duration
	err         error
	wg          sync.WaitGroup
}

var CloudronServer *CloudronServerStatus

func NewCloudronServerConnection(duration uint64) *CloudronServerStatus {
	return &CloudronServerStatus{CpuStatus: 0, MemStatus: 0, DiskStatus: 0, LastUpdated: time.Now(), err: nil, duration: time.Duration(duration), wg: sync.WaitGroup{}}
}

func main() {
	CloudronServer = NewCloudronServerConnection(5)
	StartMonitoring()
	log.Println("Starting server on 127.0.0.1:8080")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func StartMonitoring() {
	CloudronServer.wg.Add(3)

	CloudronServer.getCpu()
	CloudronServer.getMem()
	CloudronServer.getDisk()

	CloudronServer.wg.Wait()
}

func (serverStatus *CloudronServerStatus) getCpu() {
	defer serverStatus.wg.Done()
	totalPercent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		log.Println("Error getting CPU: ", err)
	}

	serverStatus.CpuStatus = totalPercent[0]
}

func (serverStatus *CloudronServerStatus) getMem() {
	defer serverStatus.wg.Done()
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}
	serverStatus.MemStatus = memInfo.UsedPercent
}

func (serverStatus *CloudronServerStatus) getDisk() {
	defer serverStatus.wg.Done()
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
	}
	serverStatus.DiskStatus = diskInfo.UsedPercent
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Cloudron stats:\n CPU Usage: %f\n Disk Usage: %f\n Memory Usage: %f\n", CloudronServer.CpuStatus, CloudronServer.DiskStatus, CloudronServer.MemStatus)

}
