package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/lua"
)

// Script manages the Lua source code for a job as well as its save locations
type Script struct {
	path    string
	source  string
	args    map[string]interface{}
	enabled bool
}

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
