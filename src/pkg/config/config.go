package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Resource struct {
	Kind     string
	Name     string
	Target   int64
	Interval int64
	Duration int64
	Slack    bool
}

type Namespace struct {
	Name     string
	Resource []Resource
}

type Config struct {
	Namespace []Namespace
}

func ReadConfig(path string) (Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if path != "" {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("Can't read config file: %s \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("config file format error: %s \n", err)
	}

	return config, nil
}
