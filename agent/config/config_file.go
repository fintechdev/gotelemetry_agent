package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Job TODO
type Job struct {
	ID         string      `toml:"id"`
	Tag        string      `toml:"tag"`
	ChannelTag string      `toml:"channel_tag"`
	Batch      bool        `toml:"batch"`
	Exec       string      `toml:"exec"`
	Script     string      `toml:"script"`
	Args       interface{} `toml:"args"`
	Template   interface{} `toml:"template"`
	Variant    string      `toml:"variant"`
	Expiration interface{} `toml:"expiration"`
	Interval   string      `toml:"interval"`
}

// ServerConfig TODO
type ServerConfig struct {
	APIToken              string      `toml:"api_token"`
	RawSubmissionInterval interface{} `toml:"submission_interval"`
}

// DataConfig TODO
type DataConfig struct {
	DataLocation *string `toml:"path"`
	TTL          *string `toml:"ttl"`
	Listen       *string `toml:"listen"`
}

// GraphiteConfig TODO
type GraphiteConfig struct {
	TCPListenPort string `toml:"listen_tcp"`
	UDPListenPort string `toml:"listen_udp"`
}

// OAuthConfigEntry TODO
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

// Interface TODO
type Interface interface {
	APIURL() string
	APIToken() (string, error)
	ChannelTag() string
	DataConfig() DataConfig
	GraphiteConfig() GraphiteConfig
	SubmissionInterval() time.Duration
	OAuthConfig() map[string]OAuthConfigEntry
	Jobs() []Job
}

// File TODO
type File struct {
	Server    ServerConfig                `toml:"server"`
	Graphite  GraphiteConfig              `toml:"graphite"`
	Data      DataConfig                  `toml:"data"`
	Listen    string                      `toml:"listen"`
	JobsField []Job                       `toml:"jobs"`
	FlowField []Job                       `toml:"flow"`
	OAuth     map[string]OAuthConfigEntry `toml:"oauth"`
}

var _ Interface = &File{}

// NewConfigFile TODO
func NewConfigFile() (*File, error) {
	source, err := ioutil.ReadFile(CLIConfig.ConfigFileLocation)

	if err != nil {
		if CLIConfig.IsPiping || CLIConfig.IsNotifying {
			return &File{
				Data:      DataConfig{},
				Graphite:  GraphiteConfig{},
				JobsField: []Job{},
			}, nil
		}

		return nil, fmt.Errorf("Unable to open configuration file at %s. Did you use --config to specify the right path?\n\n", CLIConfig.ConfigFileLocation)
	}

	result := &File{}

	_, err = toml.Decode(string(source), result)

	for _, job := range result.FlowField {
		result.JobsField = append(result.JobsField, job)
	}

	return result, err
}

// APIToken TODO
func (c *File) APIToken() (string, error) {
	result := c.Server.APIToken

	if result == "" {
		result = os.ExpandEnv("$TELEMETRY_API_TOKEN")
	}

	if result != "" {
		return result, nil
	}

	return "", errors.New("No API Token found in the configuration file or in the TELEMETRY_API_TOKEN environment variable.")
}

// APIURL TODO
func (c *File) APIURL() string {
	return CLIConfig.APIURL
}

// ChannelTag TODO
func (c *File) ChannelTag() string {
	return CLIConfig.ChannelTag
}

// DataConfig TODO
func (c *File) DataConfig() DataConfig {
	return c.Data
}

// GraphiteConfig TODO
func (c *File) GraphiteConfig() GraphiteConfig {
	return c.Graphite
}

// SubmissionInterval TODO
func (c *File) SubmissionInterval() time.Duration {
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

// Jobs TODO
func (c *File) Jobs() []Job {
	return c.JobsField
}

// OAuthConfig TODO
func (c *File) OAuthConfig() map[string]OAuthConfigEntry {
	return c.OAuth
}

// MapTemplate TODO
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
