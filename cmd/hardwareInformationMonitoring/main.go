package main

import (
	"fmt"
	"log"
	"net/http"
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
}

var CloudronServer *CloudronServerStatus

func NewCloudronServerConnection(duration uint64) *CloudronServerStatus {
	return &CloudronServerStatus{CpuStatus: 0, MemStatus: 0, DiskStatus: 0, LastUpdated: time.Now(), err: nil, duration: time.Duration(duration)}
}

func main() {
	StartMonitoring()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func StartMonitoring() {
	CloudronServer = NewCloudronServerConnection(5)
	CloudronServer.getCpu()
	CloudronServer.getMem()
	CloudronServer.getDisk()

	log.Println("CPU status:", CloudronServer.CpuStatus)
	log.Println("Memory status:", CloudronServer.MemStatus)
	log.Println("Disk status:", CloudronServer.DiskStatus)

}

func (serverStatus *CloudronServerStatus) getCpu() {
	totalPercent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		log.Println("Error getting CPU: ", err)
	}

	serverStatus.CpuStatus = totalPercent[0]
}

func (serverStatus *CloudronServerStatus) getMem() {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}
	serverStatus.MemStatus = memInfo.UsedPercent
}

func (serverStatus *CloudronServerStatus) getDisk() {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
	}
	serverStatus.DiskStatus = diskInfo.UsedPercent
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
