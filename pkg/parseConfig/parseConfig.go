package parseConfig

import (
	"errors"
	"io/ioutil"

	"log"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	"github.com/Setom29/CloudronMonitoring/pkg/report"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Args   monitor.Args  `yaml:"args"`
	Report report.Report `yaml:"report"`
}

func Parse(cfgFile string) (monitor.Args, report.Report, error) {
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Println("Error reading config file: ", err)
		return monitor.Args{}, report.Report{}, err
	}

	var cfg Config

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return monitor.Args{}, report.Report{}, err
	}

	if err = validateArgs(cfg); err != nil {
		return monitor.Args{}, report.Report{}, err
	}

	cfg.Report.Message = ""
	cfg.Report.URL = ""
	cfg.Args.CheckTime = 5
	cfg.Args.DBFile = "./data/sqlite.db"
	return cfg.Args,
		cfg.Report, nil
}

func validateArgs(cfg Config) error {
	if cfg.Args.CPULimit > 100 || cfg.Args.CPULimit < 0 {
		return errors.New("wrong value for CPU limit")
	}
	if cfg.Args.RAMLimit > 100 || cfg.Args.RAMLimit < 0 {
		return errors.New("wrong value for RAM limit")
	}
	if cfg.Args.CPULimit > 100 || cfg.Args.CPULimit < 0 {
		return errors.New("wrong value for Disk limit")
	}
	if cfg.Args.Duration < 0 {
		return errors.New("wrong value for duration")
	}
	// Parse time value
	if cfg.Args.DBClearTime > 23 || cfg.Args.DBClearTime < 0 {
		return errors.New("wrong value for db_clear_time")
	}
	if cfg.Args.MonitorLogLevel > 6 || cfg.Args.MonitorLogLevel < 1 {
		return errors.New("wrong value for monitor_log_level")
	}

	return nil
}
