package monitor

import (
	"log"
	"os"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

func StartMonitoring(f Flags, dbFile string, ch chan os.Signal) {
	db, err := dbConn.NewDB(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ch:
			db.Close()
		default:
			stat := dbConn.ServerStatus{CPUStatus: getServerStats.GetCpu(f.duration), RAMStatus: getServerStats.GetMem(), DiskStatus: getServerStats.GetDisk()}
			stat.Time = time.Now()
			if err != nil {
				log.Println(err)
			}

		}
	}
}
