package domain

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/labstack/echo"
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

func NewConfig() *Config {
	return &Config{}
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

/*
 * group.GET("", config.Get)
 */
func (c *Config) SetHandler(group *echo.Group) {
	group.GET("", c.Get)
}

// Get a config
func (c *Config) Get(context echo.Context) error {
	return nil
}
