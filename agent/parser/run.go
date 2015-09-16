package parser

import (
	"errors"
	"fmt"
)

func Run(notificationProvider executionContextNotificationProvider, jobSpawner executionContextJobSpawner, args map[string]interface{}, commands []Command) (map[string]interface{}, error) {
	ec := newExecutionContext(notificationProvider, jobSpawner, args)

	for _, cmd := range commands {
		if err := cmd.execute(ec); err != nil {
			return nil, errors.New(fmt.Sprintf("Runtime error: %s", err))
		}
	}

	return ec.output, nil
}
