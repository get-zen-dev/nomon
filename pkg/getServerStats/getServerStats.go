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
	SwapStatus       float64
	SwapTotal        uint64
	DiskStatus       float64
	DiskTotal        uint64
	duration         time.Duration
	CpuUsageDuration time.Duration
	err              error
	wg               sync.WaitGroup
	limit            int
}

func NewServerConnection(limit int, duration time.Duration) *ServerStatus {

	return &ServerStatus{
		CpuStatus:        0,
		MemStatus:        0,
		MemTotal:         0,
		SwapStatus:       0,
		SwapTotal:        0,
		DiskStatus:       0,
		DiskTotal:        0,
		duration:         duration,
		CpuUsageDuration: time.Duration(2),
		err:              nil,
		wg:               sync.WaitGroup{},
		limit:            limit}
}

func StartMonitoring(limit int, duration time.Duration) {

	serverStat := NewServerConnection(limit, duration)
	serverStat.getTotalMetrics()
	log.Printf("\nClourdon total stats:\nMemory total: %fGB\nDisk total: %fGB\nSwap total: %fGB",
		float64(serverStat.MemTotal)/1024/1024/1024,
		float64(serverStat.DiskTotal)/1024/1024/1024,
		float64(serverStat.SwapTotal)/1024/1024/1024)
	for {
		serverStat.wg.Add(3)

		go serverStat.getCpu()
		go serverStat.getMem()
		go serverStat.getDisk()

		serverStat.wg.Wait()
		log.Printf("Cloudron stats:\n CPU Usage: %f%%\n Disk Usage: %f%%\n Memory Usage: %f%%\n Swap Usage: %f%%\n",
			serverStat.CpuStatus, serverStat.DiskStatus, serverStat.MemStatus, serverStat.SwapStatus)
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
		log.Println("Error getting mem.VirtualMemory(): ", err)
	}
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		log.Println("Error getting mem.SwapMemory():", err)
	}
	serverStatus.SwapStatus = swapInfo.UsedPercent
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

	serverStatus.SwapTotal = memInfo.SwapTotal
}
