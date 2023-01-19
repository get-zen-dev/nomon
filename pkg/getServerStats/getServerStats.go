package getServerStats

import (
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type ServerStatus struct {
	CpuStatus        float64
	MemStatus        float64
	MemTotal         uint64
	DiskStatus       float64
	DiskTotal        uint64
	LastUpdated      time.Time
	duration         time.Duration
	CpuUsageDuration time.Duration
	err              error
	wg               sync.WaitGroup
}

func NewCloudronServerConnection(duration uint64) *ServerStatus {

	return &ServerStatus{
		CpuStatus:        0,
		MemStatus:        0,
		MemTotal:         0,
		DiskStatus:       0,
		DiskTotal:        0,
		LastUpdated:      time.Now(),
		duration:         time.Duration(duration),
		CpuUsageDuration: time.Duration(2),
		err:              nil,
		wg:               sync.WaitGroup{}}
}

func StartMonitoring() {

	serverStat := NewCloudronServerConnection(5)
	serverStat.getTotalMetrics()
	log.Printf("\nClourdon total stats:\nMemory total: %d\nDisk total: %d\n", serverStat.MemTotal, serverStat.DiskTotal)
	for {
		serverStat.wg.Add(3)

		go serverStat.getCpu()
		go serverStat.getMem()
		go serverStat.getDisk()

		serverStat.wg.Wait()
		serverStat.LastUpdated = time.Now()
		log.Printf("Cloudron stats:\n CPU Usage: %f\n Disk Usage: %f\n Memory Usage: %f\n", serverStat.CpuStatus, serverStat.DiskStatus, serverStat.MemStatus)
		time.Sleep((serverStat.duration - serverStat.CpuUsageDuration) * time.Second)
	}

}

func (serverStatus *ServerStatus) getCpu() {
	defer serverStatus.wg.Done()
	totalPercent, err := cpu.Percent(serverStatus.CpuUsageDuration*time.Second, false)
	if err != nil {
		log.Println("Error getting CPU: ", err)
	}

	serverStatus.CpuStatus = totalPercent[0]
}

func (serverStatus *ServerStatus) getMem() {
	defer serverStatus.wg.Done()
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}
	serverStatus.MemStatus = memInfo.UsedPercent
}

func (serverStatus *ServerStatus) getDisk() {
	defer serverStatus.wg.Done()
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
	}
	serverStatus.DiskStatus = diskInfo.UsedPercent
}
func (serverStatus *ServerStatus) getTotalMetrics() {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
	}
	serverStatus.MemTotal = memInfo.Total

	serverStatus.DiskTotal = diskInfo.Total
}
