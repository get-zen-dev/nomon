package getServerStats

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// GetCpu returns Cpu usage in percentage
func GetCpu(duration int) float64 {
	totalPercent, err := cpu.Percent(time.Duration(duration)*time.Second, false)
	if err != nil {
		log.Println("Error getting CPU: ", err)
		return 0
	}
	return totalPercent[0]
}

// GetMem returns Ram usage in percentage
func GetMem() uint64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting mem.VirtualMemory(): ", err)
		return 0
	}
	return memInfo.Used
}

// GetSwap returns Swap usage in percentage
func GetSwap() uint64 {
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		log.Println("Error getting mem.SwapMemory():", err)
		return 0
	}
	return swapInfo.Used
}

// GetDisk returns Disk usage in percentage
func GetDisk() uint64 {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
		return 0
	}
	return diskInfo.Used
}

// GetTotalMetrics returns Ram, Swap, and Disk total values
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
