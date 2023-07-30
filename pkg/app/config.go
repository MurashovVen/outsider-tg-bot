package app

import (
	"github.com/kelseyhightower/envconfig"
)

var configPrefix = "APP"

func MustProcessConfig(cfg interface{}) {
	if err := ProcessConfig(cfg); err != nil {
		panic("reading configuration: " + err.Error())
	}
}

func ProcessConfig(cfg interface{}) error {
	if err := envconfig.Usage(configPrefix, cfg); err != nil {
		return err
	}

	if err := envconfig.Process(configPrefix, cfg); err != nil {
		return err
	}

	return nil
}
