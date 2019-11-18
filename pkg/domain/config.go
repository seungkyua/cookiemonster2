package domain

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"reflect"
)

type Resource struct {
	Kind   string
	Name   string
	Target int64
}

type Namespace struct {
	Name     string
	Resource []Resource
}

type Config struct {
	Namespace []Namespace
	Interval  int64
	Duration  int64
	Slack     bool
	Change    bool
}

var config *Config
var path string

func init() {
	path := "../config"
	config = &Config{}
	if err := config.ReadConfig(path); err != nil {
		log.Println(err)
	}

}

func SlackON() bool {
	return config.Slack
}

func (c *Config) ConfigCheck(target *Config) bool {
	if reflect.DeepEqual(c, target) {
		return true
	}
	return false
}

func GetConfig() *Config {
	if err := config.ReadConfig(path); err != nil {
		log.Println(err)
	}
	return config
}

func (c *Config) ReadConfig(path string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if path != "" {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Can't read config file: %s \n", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("config file format error: %s \n", err)
	}

	return nil
}

func (c *Config) ResetConfig() {
	var configNew *Config
	configNew = &Config{}
	configNew.ReadConfig(path)
	if c.ConfigCheck(configNew) && config.Change {
		config = configNew
	}
}
