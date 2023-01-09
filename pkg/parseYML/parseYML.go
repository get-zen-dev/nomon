package parseYML

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func ParseYAML(cfgFile string) (Cfg, error) {
	file, err := ioutil.ReadFile(cfgFile)

	if err != nil {
		return Cfg{}, err
	}

	var cfg Cfg

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return Cfg{}, err
	}
	return cfg, nil
}

type Cfg struct {
	Token  string   `yaml:"Token"`
	Domain string   `yaml:"Domain"`
	Links  []string `yaml:"Links"`
}
