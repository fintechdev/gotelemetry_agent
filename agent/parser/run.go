package parser

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
)

func Run(commands []command) (map[string]interface{}, error) {
	ac, err := aggregations.GetContext()

	if err != nil {
		return nil, err
	}

	ec := newExecutionContext(ac)

	for _, cmd := range commands {
		if err := cmd.execute(ec); err != nil {
			return nil, err
		}
	}

	return ec.output, nil
}
