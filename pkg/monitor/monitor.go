package monitor

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
	"github.com/Setom29/CloudronMonitoring/pkg/report"
	log "github.com/sirupsen/logrus"
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
	CPULimit        float64 `yaml:"cpu_limit"`
	CPUCycles       int     `yaml:"cpu_cycles"`
	RAMLimit        float64 `yaml:"ram_limit"`
	RAMCycles       int     `yaml:"ram_cycles"`
	DiskLimit       float64 `yaml:"disk_limit"`
	DiskCycles      int     `yaml:"disk_cycles"`
	Duration        int     `yaml:"duration"`
	Port            int     `yaml:"port"`
	DBClearTime     int     `yaml:"db_clear_time"`
	MonitorLogLevel int     `yaml:"monitor_log_level"`
	CheckTime       int
	DBFile          string
}

func NewMonitor(args Args, r report.Report) *Monitor {
	log.Trace("monitor:NewMonitor")
	db, err := dbConn.NewDB(args.DBFile)
	if err != nil {
		log.Fatal("Error creating database: ", err)
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
	log.Debug("New monitor created")
	return &m
}

// StartMonitoring creates new db connection and pushes statistics to the database
func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	log.Trace("monitor:StartMonitoring")
	log.Info("Start monitoring")
	for {
		select {
		case s := <-ch:
			log.Debug("Got signal: ", s)
			monitor.DB.Close()
			os.Exit(0)
		default:
			log.Debug("Monitoring iteration")
			stat := dbConn.ServerStatus{
				CPUUsed:  getServerStats.GetCpu(monitor.Args.CheckTime),
				RAMUsed:  getServerStats.GetMem(),
				DiskUsed: getServerStats.GetDisk(),
				Time:     time.Now(),
			}
			err := monitor.DB.Add(stat)
			if err != nil {
				log.Error("Error adding stats: ", err)
			}
			monitor.Analyse()
			monitor.ClearDatabase()
		}
	}
}

// CheckStatus reads data from database and alerts if metric's usage limit exceeded
func (monitor *Monitor) Analyse() {
	log.Trace("monitor:Analyse")
	if time.Now().UTC().Before(monitor.LastCheck.Add(time.Duration(monitor.Args.Duration) * time.Second)) {
		return
	} else {
		monitor.LastCheck = time.Now().UTC()
	}
	rows, err := monitor.DB.Sql.Query(fmt.Sprintf("SELECT * FROM serverStatus WHERE time >= Datetime('now', '-%d seconds', 'localtime');", monitor.Args.Duration))
	if err != nil {
		log.Fatal("Error getting rows from db: ", err)
	}
	defer rows.Close()
	log.Trace("Calculating the cumulative sum for metrics")
	// cumulative sum for metrics
	counter := 0
	cpuUsedCumSum, ramUsedCumSum, diskUsedCumSum := 0.0, 0.0, 0.0
	for rows.Next() {
		counter++
		stat := dbConn.ServerStatus{}
		err := rows.Scan(&stat.Time, &stat.CPUUsed, &stat.RAMUsed, &stat.DiskUsed)
		if err != nil {
			log.Error("Error scanning rows: ", err)
			continue
		}
		cpuUsedCumSum += stat.CPUUsed
		ramUsedCumSum += float64(stat.RAMUsed)
		diskUsedCumSum += float64(stat.DiskUsed)
	}

	// alert check
	msg := monitor.AnilizeStatistics(cpuUsedCumSum/float64(counter), ramUsedCumSum/float64(counter), diskUsedCumSum/float64(counter))
	log.Debug("Alert check message: ", msg)

	if msg != "" {
		log.Warn(msg)
		monitor.Report.SendMessage(msg)
	}
}

// AnilizeStatistics gets statistics values, checks for alerts and returns message for reporting
func (monitor *Monitor) AnilizeStatistics(cpuStat, ramStat, diskStat float64) string {
	log.Trace("monitor:AnilizeStatistics")
	log.Infof("CPU: %f%%, RAM: %f%%, Disk: %f%%", cpuStat, ramStat/float64(monitor.RAMTotal)*100, diskStat/float64(monitor.DiskTotal)*100)
	msgArr := make([]string, 0)

	// CPU
	if cpuStat > monitor.Args.CPULimit {
		monitor.Counters.CPUCounter++
		if monitor.Counters.CPUCounter >= monitor.Args.CPUCycles {
			msgArr = append(msgArr, fmt.Sprintf("CPU usage limit exceeded: %f%%", cpuStat))
		}
	} else {
		if monitor.Counters.CPUCounter >= monitor.Args.CPUCycles {
			msgArr = append(msgArr, "CPU usage within normal limits.")
		}
		monitor.Counters.CPUCounter = 0
	}

	// RAM
	if ramStat/float64(monitor.RAMTotal)*100 > monitor.Args.RAMLimit {
		monitor.Counters.RAMCounter++
		if monitor.Counters.RAMCounter >= monitor.Args.RAMCycles {
			msgArr = append(msgArr, fmt.Sprintf("RAM usage limit exceeded: %f%% (%f GB /%f GB)",
				ramStat/float64(monitor.RAMTotal)*100,
				ramStat/math.Pow(1024, 3),
				float64(monitor.RAMTotal)/math.Pow(1024, 3)))
		}
	} else {
		if monitor.Counters.RAMCounter >= monitor.Args.RAMCycles {
			msgArr = append(msgArr, "RAM usage within normal limits.")
		}
		monitor.Counters.RAMCounter = 0
	}

	// Disk
	if diskStat/float64(monitor.DiskTotal)*100 > monitor.Args.DiskLimit {
		monitor.Counters.DiskCounter++
		if monitor.Counters.DiskCounter >= monitor.Args.DiskCycles {
			msgArr = append(msgArr, fmt.Sprintf("Disk usage limit exceeded: %f%% (%f GB /%f GB)",
				diskStat/float64(monitor.DiskTotal)*100,
				diskStat/math.Pow(1024, 3),
				float64(monitor.DiskTotal)/math.Pow(1024, 3)))
		}
	} else {
		if monitor.Counters.DiskCounter >= monitor.Args.DiskCycles {
			msgArr = append(msgArr, "Disk usage within normal limits.")
		}
		monitor.Counters.DiskCounter = 0
	}

	return strings.Join(msgArr, "\n")
}

func (monitor *Monitor) ClearDatabase() {
	log.Trace("monitor:ClearDatabase")
	if currTime := time.Now().UTC(); currTime.Hour() >= monitor.Args.DBClearTime && currTime.Hour() < monitor.Args.DBClearTime+1 &&
		currTime.After(monitor.LastClearTime.AddDate(0, 0, 1)) {
		log.Info("Start clearing outdated values")
		_, err := monitor.DB.Sql.Exec(fmt.Sprintf("DELETE FROM serverStatus WHERE time < Datetime('now', '-%d seconds', 'localtime');", monitor.Args.Duration))
		if err != nil {
			log.Error("Error clearing db: ", err)
		}
	}
}
