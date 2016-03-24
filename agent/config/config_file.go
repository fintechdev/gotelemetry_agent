package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"time"
)

type Job struct {
	Id         string      `toml:"id"`
	Tag        string      `toml:"tag"`
	ChannelTag string      `toml:"channel_tag"`
	Batch      bool        `toml:"batch"`
	Exec       string      `toml:"exec"`
	Script     string      `toml:"script"`
	Lua        string      `toml:"lua"`
	Url        string      `toml:"url"`
	Args       interface{} `toml:"args"`
	Template   interface{} `toml:"template"`
	Variant    string      `toml:"variant"`
	Expiration interface{} `toml:"expiration"`
	Interval   string      `toml:"interval"`
}

func (j Job) ID() string {
	if j.Id != "" {
		return j.Id
	}

	return j.Tag
}

type ServerConfig struct {
	APIToken              string      `toml:"api_token"`
	RawSubmissionInterval interface{} `toml:"submission_interval"`
}

type DataConfig struct {
	DataLocation *string `toml:"path"`
	TTL          *string `toml:"ttl"`
	Listen       *string `toml:"listen"`
}

type GraphiteConfig struct {
	TCPListenPort string `toml:"listen_tcp"`
	UDPListenPort string `toml:"listen_udp"`
}

type OAuthConfigEntry struct {
	Version          int               `toml:"version"`
	ClientID         string            `toml:"client_id"`
	ClientSecret     string            `toml:"client_secret"`
	CredentialsURL   string            `toml:"credentials_url"`
	AuthorizationURL string            `toml:"authorization_url"`
	TokenURL         string            `toml:"token_url"`
	Scopes           []string          `toml:"scopes"`
	Header           map[string]string `toml:"header"`
	SignatureMethod  string            `toml:"signature_method"`
	PrivateKey       string            `toml:"private_key"`
}

type ConfigInterface interface {
	APIURL() string
	APIToken() (string, error)
	ChannelTag() string
	DataConfig() DataConfig
	GraphiteConfig() GraphiteConfig
	SubmissionInterval() time.Duration
	OAuthConfig() map[string]OAuthConfigEntry
	Jobs() []Job
}

type ConfigFile struct {
	Server    ServerConfig                `toml:"server"`
	Graphite  GraphiteConfig              `toml:"graphite"`
	Data      DataConfig                  `toml:"data"`
	Listen    string                      `toml:"listen"`
	JobsField []Job                       `toml:"jobs"`
	FlowField []Job                       `toml:"flow"`
	OAuth     map[string]OAuthConfigEntry `toml:"oauth"`
}

var _ ConfigInterface = &ConfigFile{}

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

	_, err = toml.Decode(string(source), result)

	for _, job := range result.FlowField {
		result.JobsField = append(result.JobsField, job)
	}

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

func (c *ConfigFile) APIURL() string {
	return CLIConfig.APIURL
}

func (c *ConfigFile) ChannelTag() string {
	return CLIConfig.ChannelTag
}

func (c *ConfigFile) DataConfig() DataConfig {
	return c.Data
}

func (c *ConfigFile) GraphiteConfig() GraphiteConfig {
	return c.Graphite
}

func (c *ConfigFile) SubmissionInterval() time.Duration {
	if s, ok := c.Server.RawSubmissionInterval.(string); ok {
		d, err := ParseTimeInterval(s)

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

func (c *ConfigFile) OAuthConfig() map[string]OAuthConfigEntry {
	return c.OAuth
}

func MapTemplate(from interface{}) interface{} {
	switch from.(type) {
	case map[interface{}]interface{}:
		result := map[string]interface{}{}

		for index, value := range from.(map[interface{}]interface{}) {
			result[index.(string)] = MapTemplate(value)
		}

		return result

	case []interface{}:
		f := from.([]interface{})

		for index, value := range f {
			f[index] = MapTemplate(value)
		}

		return f

	default:
		return from
	}
}
