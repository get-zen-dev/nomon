package parseConfig

import (
	"errors"
	"io/ioutil"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	"gopkg.in/yaml.v3"
)

func ParseConfig(cfgFile string) (monitor.Args, error) {
	file, err := ioutil.ReadFile(cfgFile)

	if err != nil {
		return monitor.Args{}, err
	}

	var args monitor.Args

	err = yaml.Unmarshal(file, &args)
	if err != nil {
		return monitor.Args{}, err
	}
	if args.Limit > 100 || args.Limit < 0 {
		return monitor.Args{}, errors.New("wrong value for persentage")
	}
	if args.Duration < 0 {
		return monitor.Args{}, errors.New("wrong value for duration")
	}
	if args.CheckTime < 0 || args.CheckTime > args.Duration {
		return monitor.Args{}, errors.New("wrong value for check time")
	}
	return args, nil
}
