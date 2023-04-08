package parseConfig

import (
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	"github.com/Setom29/CloudronMonitoring/pkg/report"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Args   monitor.Args  `yaml:"args"`
	Report report.Report `yaml:"report"`
}

func Parse(cfgFile string) (*monitor.Args, *report.Report, error) {
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, nil, err
	}
	log.Println(cfg)
	if err = validateArgs(cfg); err != nil {
		return nil, nil, err
	}

	cfg.Report.Message = ""
	cfg.Report.URL = ""
	cfg.Args.CheckTime = 5
	cfg.Args.DBFile = "./data/sqlite.db"
	return &cfg.Args,
		&cfg.Report, nil
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
	_, err := time.Parse("15:04:05", cfg.Args.DBClearTime)
	if err != nil {
		return errors.New("wrong value for db_clear_time")
	}
	return nil
}
