package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/lua"
)

// Script manages the Lua source code for a job as well as its save locations
type script struct {
	filePath string
	fileMode os.FileMode
	source   string
	args     map[string]interface{}
	enabled  bool
	job      *Job
}

func newScriptFromPath(filePath string, args map[string]interface{}) (*script, error) {

	s := &script{
		enabled: true,
	}

	if !strings.HasSuffix(filePath, ".lua") {
		return nil, fmt.Errorf("Unknown script type for file `%s`", filePath)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve info for the Lua script at %s\n\n", filePath)
	}

	fileMode := info.Mode()
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	s.source = string(source)
	s.filePath = filePath
	s.fileMode = fileMode
	s.args = args

	return s, nil
}

func newScriptFromSource(source string) *script {
	s := &script{
		enabled: true,
	}

	s.source = source

	return s
}

func (s *script) UpdateExternalScript(source string) error {
	buf := new(bytes.Buffer)

	if _, err := buf.WriteString(source); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(s.filePath, buf.Bytes(), s.fileMode); err != nil {
		log.Fatal(err)
		return err
	}

	s.source = source
	return nil
}

func (s *script) exec(j *Job) (string, error) {
	output, err := lua.Exec(s.source, j, s.args)

	if err != nil {
		return "", err
	}

	out, err := json.Marshal(config.MapTemplate(output))

	return string(out), err
}
