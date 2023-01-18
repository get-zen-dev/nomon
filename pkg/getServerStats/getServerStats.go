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
	CpuStatus   float64
	MemStatus   float64
	DiskStatus  float64
	LastUpdated time.Time
	duration    time.Duration
	err         error
	wg          sync.WaitGroup
}

func NewCloudronServerConnection(duration uint64) *ServerStatus {
	return &ServerStatus{CpuStatus: 0, MemStatus: 0, DiskStatus: 0, LastUpdated: time.Now(), err: nil, duration: time.Duration(duration), wg: sync.WaitGroup{}}
}

func StartMonitoring() *ServerStatus {

	serverStat := NewCloudronServerConnection(5)
	serverStat.wg.Add(3)

	go serverStat.getCpu()
	go serverStat.getMem()
	go serverStat.getDisk()

	serverStat.wg.Wait()
	return serverStat
}

func (serverStatus *ServerStatus) getCpu() {
	defer serverStatus.wg.Done()
	totalPercent, err := cpu.Percent(3*time.Second, false)
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
