package functions

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions/schemas"
)

func init() {
	schemas.LoadSchema("extract")
	functionHandlers["$extract"] = extractHandler
}

func extractHandler(context *aggregations.Context, input interface{}) (interface{}, error) {
	if err := validatePayload("$extract", input); err != nil {
		return nil, err
	}

	data := input.(map[string]interface{})
	props := data["props"].([]string)

	defaults, hasDefault := data["default"].(map[string]interface{})

	if obj, ok := data["from"].(map[string]interface{}); ok {
		result := map[string]interface{}{}

		for _, prop := range props {
			if val, ok := obj[prop]; ok {
				result[prop] = val
			} else if hasDefault {
				if val, ok := defaults[prop]; ok {
					result[prop] = val
				} else {
					return nil, errors.New(fmt.Sprintf("Property %s not found in %#v, and no default value defined", prop, obj))
				}
			} else {
				return nil, errors.New(fmt.Sprintf("Property %s not found in %#v, and no default value defined", prop, obj))
			}
		}

		return result, nil
	}

	if obj, ok := data["from"].([]interface{}); ok {
		result := make([]interface{}, len(obj))

		for index, rec := range obj {
			if record, ok := rec.(map[string]interface{}); ok {
				out := map[string]interface{}{}

				for _, prop := range props {
					if val, ok := record[prop]; ok {
						out[prop] = val
					} else if hasDefault {
						if val, ok := defaults[prop]; ok {
							out[prop] = val
						} else {
							return nil, errors.New(fmt.Sprintf("Property %s not found in %#v at index %d, and no default value defined", prop, record, index))
						}
					} else {
						return nil, errors.New(fmt.Sprintf("Property %s not found in %#v at index %d, and no default value defined", prop, record, index))
					}

					result[index] = out
				}
			} else {
				return nil, errors.New(fmt.Sprintf("The item at index %d is not a hash.", index))
			}
		}

		return result, nil
	}

	return nil, errors.New(fmt.Sprintf("$extract doesn't know how to handle %#v", data["from"]))
}
