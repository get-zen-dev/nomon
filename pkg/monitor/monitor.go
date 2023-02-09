package monitor

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

type Monitor struct {
	DB *dbConn.DB
	F  Args
}

type Args struct {
	Limit     float64 `yaml:"limit"`
	Duration  int     `yaml:"duration"`
	CheckTime int     `yaml:"checktime"`
	DBFile    string  `yaml:"dbfile"`
	Port      string  `yaml:"port"`
}

func NewMonitor(f Args) *Monitor {
	db, err := dbConn.NewDB(f.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	m := Monitor{
		DB: db,
		F:  f,
	}
	return &m
}

// StartMonitoring creates new db connection and pushes statistics to the database
func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	log.Println("Starting monitoring")
	db, err := dbConn.NewDB(monitor.F.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ch:
			db.Close()
			return
		default:
			stat := dbConn.ServerStatus{
				CPUStatus:  getServerStats.GetCpu(monitor.F.CheckTime),
				RAMStatus:  getServerStats.GetMem(),
				DiskStatus: getServerStats.GetDisk(),
				Time:       time.Now()}
			if err != nil {
				log.Println(err)
			}
			err = monitor.DB.Add(stat)
			if err != nil {
				log.Println(err)
			}
			monitor.CheckStatus()
		}
	}
}

// CheckStatus reads data from database and alerts if metric's usage limit exceeded
func (monitor *Monitor) CheckStatus() {
	monitor.DB.Sql.Exec(fmt.Sprintf("DELETE FROM serverStatus WHERE time < Datetime('now', '-%d seconds', 'localtime');", monitor.F.Duration))
	rows, err := monitor.DB.Sql.Query(fmt.Sprintf("SELECT * FROM serverStatus WHERE time >= Datetime('now', '-%d seconds', 'localtime');", monitor.F.Duration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	counter := 0
	CPUCumSum, RAMCumSum, DiskCumSum := 0.0, 0.0, 0.0
	for rows.Next() {
		counter++
		stat := dbConn.ServerStatus{}
		err := rows.Scan(&stat.Time, &stat.CPUStatus, &stat.RAMStatus, &stat.DiskStatus)
		if err != nil {
			fmt.Println(err)
			continue
		}
		CPUCumSum += stat.CPUStatus
		RAMCumSum += stat.RAMStatus
		DiskCumSum += stat.DiskStatus
	}
	if CPUCumSum/float64(counter) > monitor.F.Limit {
		log.Println("CPU usage limit exceeded")
	}
	if RAMCumSum/float64(counter) > monitor.F.Limit {
		log.Println("RAM usage limit exceeded")
	}
	if DiskCumSum/float64(counter) > monitor.F.Limit {
		log.Println("Disk usage limit exceeded")
	}
}
