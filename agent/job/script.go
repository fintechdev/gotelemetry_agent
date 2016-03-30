package job

import (
	"encoding/json"
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/lua"
	"io/ioutil"
	"strings"
)

// Struct ProcessPlugin allows the agent to execute an external process and use its
// output as data that can be fed to the Telemetry API.
type Script struct {
	path    string
	source  string
	args    map[string]interface{}
	enabled bool
}

// The script struct contains the Lua code that is \
func newScript(path string, args map[string]interface{}) (*Script, error) {

	s := &Script{
		enabled: true,
	}

	if !strings.HasSuffix(path, ".lua") {
		return nil, fmt.Errorf("Unknown script type for file `%s`", path)
	}

	s.path = path
	s.args = args

	source, err := ioutil.ReadFile(s.path)

	if err != nil {
		return nil, err
	}

	s.source = string(source)

	return s, nil
}

func (s *Script) exec(j *Job) (string, error) {

	output, err := lua.Exec(s.source, j, s.args)

	if err != nil {
		return "", err
	}

	out, err := json.Marshal(config.MapTemplate(output))

	return string(out), err
}
