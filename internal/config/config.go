package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Endpoint    string         `yaml:"endpoint"`
	MaxPollTime time.Duration  `yaml:"maxPollTime"`
	Rotate      RotationConfig `yaml:"rotate"`
	OutputDir   string         `yaml:"outputDir"`
}

type RotationConfig struct {
	Interval time.Duration `yaml:"interval"`
	Size     int64         `yaml:"size"`
}

func New(env, runId string) (*Config, error) {
	c := &Config{
		MaxPollTime: time.Duration(60 * time.Second), // default 1 minute
		Rotate: RotationConfig{
			Interval: time.Duration(300 * time.Second), // default 5 minutes
			Size:     1024 * 1024,                      // defualt 1MB
		},
	}

	yb, err := os.ReadFile("config/" + env + ".yaml")
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yb, c); err != nil {
		return nil, err
	}
	c.OutputDir = fmt.Sprintf("%s/%s", c.OutputDir, runId)

	return c, nil
}
