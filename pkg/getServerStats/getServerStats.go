package getServerStats

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func GetCpu(duration time.Duration) float64 {
	totalPercent, err := cpu.Percent(duration, false)
	if err != nil {
		log.Println("Error getting CPU: ", err)
		return 0
	}
	return totalPercent[0]
}

func GetMem() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting mem.VirtualMemory(): ", err)
		return 0
	}
	return memInfo.UsedPercent
}

func GetSwap() float64 {
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		log.Println("Error getting mem.SwapMemory():", err)
		return 0
	}
	return swapInfo.UsedPercent
}

func GetDisk() float64 {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
		return 0
	}
	return diskInfo.UsedPercent
}

func GetTotalMetrics() (uint64, uint64, uint64) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting Memory: ", err)
		return 0, 0, 0
	}
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
		return 0, 0, 0
	}
	return memInfo.Total, memInfo.SwapTotal, diskInfo.Total
}
