package monitor

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
	"github.com/Setom29/CloudronMonitoring/pkg/report"
)

type Monitor struct {
	DB            *dbConn.DB
	Args          Args
	Report        report.Report
	Counters      Counters
	WG            *sync.WaitGroup
	RAMTotal      uint64
	DiskTotal     uint64
	LastCheck     time.Time
	LastClearTime time.Time
}
type Counters struct {
	CPUCounter  int
	RAMCounter  int
	DiskCounter int
}

type Args struct {
	CPULimit    float64 `yaml:"cpu_limit"`
	CPUCycles   int     `yaml:"cpu_cycles"`
	RAMLimit    float64 `yaml:"ram_limit"`
	RAMCycles   int     `yaml:"ram_cycles"`
	DiskLimit   float64 `yaml:"disk_limit"`
	DiskCycles  int     `yaml:"disk_cycles"`
	Duration    int     `yaml:"duration"`
	Port        int     `yaml:"port"`
	DBClearTime int     `yaml:"db_clear_time"`
	CheckTime   int
	DBFile      string
}

func NewMonitor(args Args, r report.Report) *Monitor {
	db, err := dbConn.NewDB(args.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	ramTotal, _, diskTotal := getServerStats.GetTotalMetrics()
	t := time.Now().UTC()
	m := Monitor{
		DB:       db,
		Args:     args,
		Report:   r,
		Counters: Counters{CPUCounter: 0, RAMCounter: 0, DiskCounter: 0},
		WG:       &sync.WaitGroup{},

		RAMTotal:      ramTotal,
		DiskTotal:     diskTotal,
		LastCheck:     t,
		LastClearTime: t,
	}
	return &m
}

// StartMonitoring creates new db connection and pushes statistics to the database
func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	log.Println("Start monitoring")

	for {
		select {
		case <-ch:
			monitor.DB.Close()
			os.Exit(0)
		default:

			stat := dbConn.ServerStatus{
				CPUUsed:  getServerStats.GetCpu(monitor.Args.CheckTime),
				RAMUsed:  getServerStats.GetMem(),
				DiskUsed: getServerStats.GetDisk(),
				Time:     time.Now()}
			err := monitor.DB.Add(stat)
			if err != nil {
				log.Println("Error adding stats:", err)
			}
			monitor.Analyse()
			monitor.ClearDatabase()
		}
	}
}

// CheckStatus reads data from database and alerts if metric's usage limit exceeded
func (monitor *Monitor) Analyse() {
	if time.Now().UTC().Before(monitor.LastCheck.Add(time.Duration(monitor.Args.Duration) * time.Second)) {
		return
	} else {
		monitor.LastCheck = time.Now().UTC()
	}
	rows, err := monitor.DB.Sql.Query(fmt.Sprintf("SELECT * FROM serverStatus WHERE time >= Datetime('now', '-%d seconds', 'localtime');", monitor.Args.Duration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// cumulative sum for metrics
	counter := 0
	cpuUsedCumSum, ramUsedCumSum, diskUsedCumSum := 0.0, 0.0, 0.0
	for rows.Next() {
		counter++
		stat := dbConn.ServerStatus{}
		err := rows.Scan(&stat.Time, &stat.CPUUsed, &stat.RAMUsed, &stat.DiskUsed)
		if err != nil {
			log.Println("Error scanning server status:", err)
			continue
		}
		cpuUsedCumSum += stat.CPUUsed
		ramUsedCumSum += float64(stat.RAMUsed)
		diskUsedCumSum += float64(stat.DiskUsed)
	}

	// alert check
	msg := monitor.AnilizeStatistics(cpuUsedCumSum/float64(counter), ramUsedCumSum/float64(counter), diskUsedCumSum/float64(counter))

	if msg != "" {
		log.Println(msg)
		monitor.Report.SendMessage(msg)
	}
}

// AnilizeStatistics gets statistics values, checks for alerts and returns message for reporting
func (monitor *Monitor) AnilizeStatistics(cpuStat, ramStat, diskStat float64) string {
	msgArr := make([]string, 0)
	if cpuStat > monitor.Args.CPULimit {
		if monitor.Counters.CPUCounter >= monitor.Args.CPUCycles {
			msgArr = append(msgArr, fmt.Sprintf("CPU usage limit exceeded: %f%%", cpuStat))
		}
		monitor.Counters.CPUCounter++
	} else {
		if monitor.Counters.CPUCounter > monitor.Args.CPUCycles {
			msgArr = append(msgArr, "CPU usage within normal limits.")
		}
		monitor.Counters.CPUCounter = 0
	}
	if ramStat/float64(monitor.RAMTotal)*100 > monitor.Args.RAMLimit {
		if monitor.Counters.RAMCounter >= monitor.Args.RAMCycles {
			msgArr = append(msgArr, fmt.Sprintf("RAM usage limit exceeded: %f%% (%f GB /%f GB)",
				ramStat/float64(monitor.RAMTotal)*100,
				ramStat/math.Pow(1024, 3),
				float64(monitor.RAMTotal)/math.Pow(1024, 3)))
		}
		monitor.Counters.RAMCounter++
	} else {
		if monitor.Counters.RAMCounter > monitor.Args.RAMCycles {
			msgArr = append(msgArr, "RAM usage within normal limits.")
		}
		monitor.Counters.CPUCounter = 0
	}
	if diskStat/float64(monitor.DiskTotal)*100 > monitor.Args.DiskLimit {
		if monitor.Counters.DiskCounter >= monitor.Args.DiskCycles {
			msgArr = append(msgArr, fmt.Sprintf("Disk usage limit exceeded: %f%% (%f GB /%f GB)",
				diskStat/float64(monitor.DiskTotal)*100,
				diskStat/math.Pow(1024, 3),
				float64(monitor.DiskTotal)/math.Pow(1024, 3)))
		}
		monitor.Counters.DiskCounter++
	} else {
		if monitor.Counters.RAMCounter > monitor.Args.RAMCycles {
			msgArr = append(msgArr, "Disk usage within normal limits.")
		}
		monitor.Counters.DiskCounter = 0
	}

	return strings.Join(msgArr, "\n")
}

func (monitor *Monitor) ClearDatabase() {
	if currTime := time.Now().UTC(); currTime.Hour() >= monitor.Args.DBClearTime && currTime.Hour() < monitor.Args.DBClearTime+1 &&
		currTime.After(monitor.LastClearTime.AddDate(0, 0, 1)) {
		log.Println("Start clearing outdated values")
		_, err := monitor.DB.Sql.Exec(fmt.Sprintf("DELETE FROM serverStatus WHERE time < Datetime('now', '-%d seconds', 'localtime');", monitor.Args.Duration))
		if err != nil {
			log.Println(err)
		}
	}
}
