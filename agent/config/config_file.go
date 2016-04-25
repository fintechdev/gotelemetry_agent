package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Job TODO
type Job struct {
	ID         string      `toml:"id"          json:"id"`
	Tag        string      `toml:"tag"         json:"tag"`
	ChannelTag string      `toml:"channel_tag" json:"channel_tag"`
	Batch      bool        `toml:"batch"       json:"batch"`
	Exec       string      `toml:"exec"        json:"exec"`
	Script     string      `toml:"script"      json:"script"`
	Args       interface{} `toml:"args"        json:"args"`
	Template   interface{} `toml:"template"    json:"template"`
	Variant    string      `toml:"variant"     json:"variant"`
	Expiration interface{} `toml:"expiration"  json:"expiration"`
	Interval   string      `toml:"interval"    json:"interval"`
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
	TTL              string            `toml:"ttl"`
}

// ServerConfig TODO
type ServerConfig struct {
	APIToken              string      `toml:"api_token"`
	RawSubmissionInterval interface{} `toml:"submission_interval"`
}

// DataConfig TODO
type DataConfig struct {
	DataLocation string `toml:"path"`
	TTL          string `toml:"ttl"`
}

// ListenerConfig TODO
type ListenerConfig struct {
	Listen    string `toml:"listen"`
	AuthToken string `toml:"auth_token"`
}

// GraphiteConfig TODO
type GraphiteConfig struct {
	TCPListenPort string `toml:"listen_tcp"`
	UDPListenPort string `toml:"listen_udp"`
}

// Interface TODO
type Interface interface {
	APIURL() string
	APIToken() string
	SetAPIToken(string)
	SetAuthToken(string)
	SetListen(string)
	SetDatabaseTTL(string)
	ChannelTag() string
	DataConfig() DataConfig
	GraphiteConfig() GraphiteConfig
	SubmissionInterval() time.Duration
	OAuthConfig() map[string]OAuthConfigEntry
	Jobs() []Job
	Listen() string
	AuthToken() string
	DatabasePath() string
	DatabaseTTL() string
}

// File TODO
type File struct {
	filePath  string
	fileMode  os.FileMode
	Listener  ListenerConfig              `toml:"listener"`
	Server    ServerConfig                `toml:"server"`
	Graphite  GraphiteConfig              `toml:"graphite"`
	Data      DataConfig                  `toml:"data"`
	JobsField []Job                       `toml:"jobs"`
	FlowField []Job                       `toml:"flow"`
	OAuth     map[string]OAuthConfigEntry `toml:"oauth"`
}

var _ Interface = &File{}

// NewConfigFile TODO
func NewConfigFile() (*File, error) {
	filePath := CLIConfig.ConfigFileLocation

	result := &File{}

	if len(filePath) == 0 {
		// Return an empty configfile if a path has not been set
		return result, nil
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to get configuration file info at %s. Did you use --config to specify the right path?\n\n", filePath)
	}

	fileMode := info.Mode()

	source, err := ioutil.ReadFile(filePath)

	if err != nil {
		if CLIConfig.IsPiping || CLIConfig.IsNotifying {
			return &File{
				Data:      DataConfig{},
				Graphite:  GraphiteConfig{},
				JobsField: []Job{},
			}, nil
		}

		return nil, fmt.Errorf("Unable to open configuration file at %s. Did you use --config to specify the right path?\n\n", filePath)
	}

	_, err = toml.Decode(string(source), result)

	for _, job := range result.FlowField {
		result.JobsField = append(result.JobsField, job)
	}

	result.filePath = filePath
	result.fileMode = fileMode

	return result, err
}

// APIToken TODO
func (c *File) APIToken() string {
	result := c.Server.APIToken

	if result == "" {
		return os.ExpandEnv("$TELEMETRY_API_TOKEN")
	}

	return result
}

// SetAPIToken replaces the config API token
func (c *File) SetAPIToken(token string) {
	c.Server.APIToken = token
}

// SetAuthToken replaces the config auth token
func (c *File) SetAuthToken(token string) {
	c.Listener.AuthToken = token
}

// SetListen TODO
func (c *File) SetListen(portNumber string) {
	c.Listener.Listen = portNumber
}

// SetDatabaseTTL TODO
func (c *File) SetDatabaseTTL(ttlString string) {
	c.Data.TTL = ttlString
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

// Listen TODO
func (c *File) Listen() string {
	if cliListenPort := CLIConfig.AuthenticationPort; len(cliListenPort) > 0 {
		return cliListenPort
	}

	return c.Listener.Listen
}

// AuthToken TODO
func (c *File) AuthToken() string {
	if cliAuthKey := CLIConfig.AuthenticationToken; len(cliAuthKey) > 0 {
		return cliAuthKey
	}

	return c.Listener.AuthToken
}

// DatabasePath TODO
func (c *File) DatabasePath() string {
	if databasePath := CLIConfig.DatabasePath; len(databasePath) > 0 {
		return databasePath
	}

	if len(c.Data.DataLocation) == 0 {
		return "agent.db"
	}

	return c.Data.DataLocation
}

// DatabaseTTL TODO
func (c *File) DatabaseTTL() string {
	if dataTTL := CLIConfig.DatabaseTTL; len(dataTTL) > 0 {
		return dataTTL
	}

	return c.Data.TTL
}
