package cloudronStatus

import (
	"time"
)

type CloudronStatusInfo struct {
	Domain      string
	CpuStatus   int
	RamStatus   int
	DiskStatus  int
	LastUpdated time.Time
}

func NewCloudronStatusInfo(domain string) *CloudronStatusInfo {

	return &CloudronStatusInfo{Domain: domain, CpuStatus: 0, RamStatus: 0, DiskStatus: 0, LastUpdated: time.Now()}
}
