package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigFile struct {
	Data        DataConfig
	Graphite    GraphiteConfig
	AllAccounts []AccountConfig
}

func NewConfigFile() (*ConfigFile, error) {
	source, err := ioutil.ReadFile(CLIConfig.ConfigFileLocation)

	if err != nil {
		if CLIConfig.IsPiping || CLIConfig.IsNotifying {
			return &ConfigFile{
				Data:        DataConfig{},
				Graphite:    GraphiteConfig{},
				AllAccounts: []AccountConfig{AccountConfig{}},
			}, nil
		}

		return nil, errors.New(fmt.Sprintf("Unable to open configuration file at %s. Did you use --config to specify the right path?\n\n", CLIConfig.ConfigFileLocation))
	}

	result := &AccountConfig{}

	err = yaml.Unmarshal(source, result)

	if err != nil {
		return nil, err
	}

	result.Jobs = make([]Job, len(result.RawJobs))

	for index, rawJob := range result.RawJobs {
		if _, ok := rawJob["config"].(map[interface{}]interface{}); ok {
			if b, err := yaml.Marshal(rawJob); err == nil {
				j := Job{}

				if err := yaml.Unmarshal(b, &j); err == nil {
					result.Jobs[index] = j
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			result.Jobs[index] = Job{Config: rawJob}
		}
	}

	return &ConfigFile{
		Data:        result.Data,
		Graphite:    result.Graphite,
		AllAccounts: []AccountConfig{*result},
	}, err
}

func (c *ConfigFile) Accounts() []AccountConfig {
	return c.AllAccounts
}

func (c *ConfigFile) DataConfig() DataConfig {
	return c.Data
}

func (c *ConfigFile) GraphiteConfig() GraphiteConfig {
	return c.Graphite
}

func (c *ConfigFile) ListenAddress() string {
	return c.AllAccounts[0].Listen
}
