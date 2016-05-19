package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// File is the top level struct of the Agent config file
type File struct {
	Listener  ListenerConfig              `toml:"listener"`
	Server    ServerConfig                `toml:"server"`
	Graphite  GraphiteConfig              `toml:"graphite"`
	Data      DataConfig                  `toml:"data"`
	JobsField []Job                       `toml:"jobs"`
	FlowField []Job                       `toml:"flow"`
	OAuth     map[string]OAuthConfigEntry `toml:"oauth"`
}

// Job handles all job and flow parameters
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

// OAuthConfigEntry handles OAuth configuration parameters
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

// ServerConfig handles configuration info for the API connection to TelemetryTV
type ServerConfig struct {
	APIToken              string      `toml:"api_token"`
	RawSubmissionInterval interface{} `toml:"submission_interval"`
}

// DataConfig handles the configuration info for the Agent's internal database
type DataConfig struct {
	DataLocation string `toml:"path"`
	TTL          string `toml:"ttl"`
}

// ListenerConfig handles configuration info for the Agent's internal API
type ListenerConfig struct {
	Listen   string `toml:"listen"`
	AuthKey  string `toml:"auth_key"`
	CertFile string `toml:"certfile"`
	KeyFile  string `toml:"keyfile"`
}

// GraphiteConfig handles the port numbers for the Agent's Graphite interface
type GraphiteConfig struct {
	TCPListenPort string `toml:"listen_tcp"`
	UDPListenPort string `toml:"listen_udp"`
}

// Interface limits the configuration scope to include only the get/set functions
type Interface interface {
	APIURL() string
	APIToken() string
	SetAPIToken(string)
	SetAuthKey(string)
	SetListen(string)
	SetDatabaseTTL(string)
	SetUDPListenPort(string)
	SetTCPListenPort(string)
	ChannelTag() string
	DataConfig() DataConfig
	GraphiteConfig() GraphiteConfig
	SubmissionInterval() time.Duration
	OAuthConfig() map[string]OAuthConfigEntry
	Jobs() []Job
	Listen() string
	AuthKey() string
	DatabasePath() string
	DatabaseTTL() string
	CertFile() string
	KeyFile() string
}

var _ Interface = &File{}

// NewConfigFile initializes a configFile object based on a configuration file path.
// returns the populated object if the parsing is successful. Returns an empty object
// if no path is provided
func NewConfigFile() (*File, error) {
	filePath := CLIConfig.ConfigFileLocation

	result := &File{}

	if len(filePath) == 0 {
		// Return an empty config file if a path has not been set
		return result, nil
	}

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

	return result, err
}

// MapTemplate takes a map of default flow values and produces an interface for
// gotelemetry to create those flows if they do not exist
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

// SetAPIToken unconditionally replaces the config API token
func (c *File) SetAPIToken(token string) {
	c.Server.APIToken = token
}

// SetAuthKey unconditionally replaces the config auth key
func (c *File) SetAuthKey(key string) {
	c.Listener.AuthKey = key
}

// SetListen unconditionally replaces the server listening port
func (c *File) SetListen(portNumber string) {
	c.Listener.Listen = portNumber
}

// SetDatabaseTTL unconditionally replaces the database series TTL
func (c *File) SetDatabaseTTL(ttlString string) {
	c.Data.TTL = ttlString
}

// SetUDPListenPort unconditionally replaces the graphite UDP port number
func (c *File) SetUDPListenPort(portString string) {
	c.Graphite.UDPListenPort = portString
}

// SetTCPListenPort unconditionally replaces the graphite TCP port number
func (c *File) SetTCPListenPort(portString string) {
	c.Graphite.TCPListenPort = portString
}

// APIToken returns the APIToken value from the configFile if present.
// Returns the $TELEMETRY_API_TOKEN environment variable value otherwise
func (c *File) APIToken() string {
	result := c.Server.APIToken

	if len(result) == 0 {
		return os.ExpandEnv("$TELEMETRY_API_TOKEN")
	}

	return result
}

// APIURL returns the APIURL value form the command line parameters
func (c *File) APIURL() string {
	return CLIConfig.APIURL
}

// ChannelTag returns the ChannelTag value from the command line parameters
func (c *File) ChannelTag() string {
	return CLIConfig.ChannelTag
}

// DataConfig returns the Data object from the configFile
func (c *File) DataConfig() DataConfig {
	return c.Data
}

// GraphiteConfig returns the Graphite object from the configFile
func (c *File) GraphiteConfig() GraphiteConfig {
	return c.Graphite
}

// SubmissionInterval parses the raw interval time from the configFile and
// returns the interval in duration format
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

// Jobs returns an array of JobsField objects from the configFile
func (c *File) Jobs() []Job {
	return c.JobsField
}

// OAuthConfig returns an array of the OAuth objects from the configFile
func (c *File) OAuthConfig() map[string]OAuthConfigEntry {
	return c.OAuth
}

// Listen returns the Listen value from the command line parameters if present and from
// the configFile object otherwise
func (c *File) Listen() string {
	if cliListenPort := CLIConfig.AuthenticationPort; len(cliListenPort) > 0 {
		return cliListenPort
	}

	return c.Listener.Listen
}

// AuthKey returns the AuthKey value from the command line parameters if present and from
// the configFile object otherwise
func (c *File) AuthKey() string {
	if cliAuthKey := CLIConfig.AuthenticationKey; len(cliAuthKey) > 0 {
		return cliAuthKey
	}

	return c.Listener.AuthKey
}

// DatabasePath returns the DatabasePath value from the command line parameters if present and from
// the configFile object otherwise. Sets a default value of agent.db in the current directory if nil
func (c *File) DatabasePath() string {
	if databasePath := CLIConfig.DatabasePath; len(databasePath) > 0 {
		return databasePath
	}

	if len(c.Data.DataLocation) == 0 {
		return "agent.db"
	}

	return c.Data.DataLocation
}

// DatabaseTTL returns the DatabaseTTL value from the command line parameters if present and from
// the configFile object otherwise
func (c *File) DatabaseTTL() string {
	if dataTTL := CLIConfig.DatabaseTTL; len(dataTTL) > 0 {
		return dataTTL
	}

	return c.Data.TTL
}

// CertFile returns the CertFile value from the command line parameters if present and from
// the configFile object otherwise
func (c *File) CertFile() string {
	if certFile := CLIConfig.CertFile; len(certFile) > 0 {
		return certFile
	}

	return c.Listener.CertFile
}

// KeyFile returns the KeyFile value from the command line parameters if present and from
// the configFile object otherwise
func (c *File) KeyFile() string {
	if keyFile := CLIConfig.KeyFile; len(keyFile) > 0 {
		return keyFile
	}

	return c.Listener.KeyFile
}
