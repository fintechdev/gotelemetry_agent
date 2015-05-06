package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type Job map[string]interface{}

func (j Job) ID() string {
	if result, success := j["id"].(string); success {
		return result
	}

	if result, success := j["flow_tag"].(string); success {
		return result
	}

	return ""
}

func (j Job) Plugin() string {
	if result, success := j["plugin"].(string); success {
		return result
	}

	return "com.telemetryapp.process"
}

type ServerConfig struct {
	APIToken              string      `yaml:"api_token"`
	RawSubmissionInterval interface{} `yaml:"submission_interval"`
}

type DataConfig struct {
	DataLocation *string `yaml:"path"`
}

type GraphiteConfig struct {
	TCPListenPort string `yaml:"listen_tcp"`
	UDPListenPort string `yaml:"listen_udp"`
}

type ConfigInterface interface {
	APIToken() (string, error)
	ListenAddress() string
	DataConfig() DataConfig
	GraphiteConfig() GraphiteConfig
	SubmissionInterval() time.Duration
	Jobs() []Job
}

type ConfigFile struct {
	Server    ServerConfig   `yaml:"server"`
	Graphite  GraphiteConfig `yaml:"graphite"`
	Data      DataConfig     `yaml:"data"`
	Listen    string         `yaml:"listen"`
	JobsField []Job          `yaml:"jobs"`
}

func NewConfigFile() (*ConfigFile, error) {
	source, err := ioutil.ReadFile(CLIConfig.ConfigFileLocation)

	if err != nil {
		if CLIConfig.IsPiping || CLIConfig.IsNotifying {
			return &ConfigFile{
				Data:      DataConfig{},
				Graphite:  GraphiteConfig{},
				JobsField: []Job{},
			}, nil
		}

		return nil, errors.New(fmt.Sprintf("Unable to open configuration file at %s. Did you use --config to specify the right path?\n\n", CLIConfig.ConfigFileLocation))
	}

	result := &ConfigFile{}

	err = yaml.Unmarshal(source, result)

	return result, err
}

func (c *ConfigFile) APIToken() (string, error) {
	result := c.Server.APIToken

	if result == "" {
		result = os.ExpandEnv("$TELEMETRY_API_TOKEN")
	}

	if result != "" {
		return result, nil
	}

	return "", errors.New("No API Token found in the configuration file or in the TELEMETRY_API_TOKEN environment variable.")
}

func (c *ConfigFile) DataConfig() DataConfig {
	return c.Data
}

func (c *ConfigFile) GraphiteConfig() GraphiteConfig {
	return c.Graphite
}

func (c *ConfigFile) ListenAddress() string {
	return c.Listen
}

func (c *ConfigFile) SubmissionInterval() time.Duration {
	if s, ok := c.Server.RawSubmissionInterval.(string); ok {
		d, err := time.ParseDuration(s)

		if err == nil {
			return d
		}
	}

	if s, ok := c.Server.RawSubmissionInterval.(float64); ok {
		return time.Duration(s) * time.Second
	}

	return 0
}

func (c *ConfigFile) Jobs() []Job {
	return c.JobsField
}
