package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Number

type globalExpression struct {
	l int
	p int
}

func newGlobalExpression(line, position int) expression {
	result := &globalExpression{line, position}

	return result
}

func (g *globalExpression) evaluate(c *executionContext) (interface{}, error) {
	return g, nil
}

// Properties

type globalProperty func(g *globalExpression) expression

var globalProperties = map[string]globalProperty{
	"counter": func(g *globalExpression) expression {
		return g.counter()
	},
	"now": func(g *globalExpression) expression {
		return g.now()
	},
	"notify": func(g *globalExpression) expression {
		return g.notify()
	},
	"series": func(g *globalExpression) expression {
		return g.series()
	},
	"anomaly": func(g *globalExpression) expression {
		return g.anomaly()
	},
	"arg": func(g *globalExpression) expression {
		return g.arg()
	},
	"load": func(g *globalExpression) expression {
		return g.load()
	},
	"spawn": func(g *globalExpression) expression {
		return g.spawn()
	},
	"log": func(g *globalExpression) expression {
		return g.log()
	},
	"error": func(g *globalExpression) expression {
		return g.emitError()
	},
	"excel": func(g *globalExpression) expression {
		return g.excel()
	},
	"get": func(g *globalExpression) expression {
		return g.get()
	},
	"post": func(g *globalExpression) expression {
		return g.post()
	},
}

func (g *globalExpression) now() expression {
	return newCallableExpression(
		"now",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return newNumericExpression(time.Now().Unix(), g.l, g.p), nil
		},
		map[string]callableArgument{},
		g.l,
		g.p,
	)
}

func (g *globalExpression) emitError() expression {
	return newCallableExpression(
		"error",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return nil, errors.New(args["message"].(string))
		},
		map[string]callableArgument{
			"message": callableArgumentString,
		},
		g.l,
		g.p,
	)
}

func (g *globalExpression) log() expression {
	return newCallableExpression(
		"log",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			log.Println(args["message"].(string))

			return nil, nil
		},
		map[string]callableArgument{
			"message": callableArgumentString,
		},
		g.l,
		g.p,
	)
}

func (g *globalExpression) excel() expression {
	return newCallableExpression(
		"excel",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			path := args["path"].(string)

			return newExcelExpression(path, g.l, g.p), nil
		},
		map[string]callableArgument{
			"path": callableArgumentString,
		},
		g.l,
		g.p,
	)
}

func (g *globalExpression) counter() expression {
	return newCallableExpression(
		"counter",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			n := args["name"].(string)
			res, _, err := aggregations.GetCounter(c.aggregationContext, n)

			if err != nil {
				return nil, err
			}

			return newCounterExpression(n, res, g.p, g.p), nil
		},
		map[string]callableArgument{"name": callableArgumentString},
		g.l,
		g.p,
	)
}

func (g *globalExpression) arg() expression {
	return newCallableExpression(
		"arg",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			n := args["name"].(string)

			if res, ok := c.arguments[n]; ok {
				return expressionFromInterface(res, g.l, g.p)
			} else {
				return nil, errors.New(fmt.Sprintf("Unknown argument %s", n))
			}
		},
		map[string]callableArgument{"name": callableArgumentString},
		g.l,
		g.p,
	)
}

func (g *globalExpression) series() expression {
	return newCallableExpression(
		"series",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			n := args["name"].(string)
			s, _, err := aggregations.GetSeries(c.aggregationContext, n)

			if err != nil {
				return nil, err
			}

			return newSeriesExpression(n, s, g.p, g.p), nil
		},
		map[string]callableArgument{"name": callableArgumentString},
		g.l,
		g.p,
	)
}

func (g *globalExpression) notify() expression {
	return newCallableExpression(
		"notify",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			title := args["title"].(string)
			message := args["message"].(string)

			dString, _ := args["duration"].(string)
			if len(dString) == 0 {
				dString = "1s"
			}
			d, err := time.ParseDuration(dString)
			if err != nil {
				return nil, err
			}

			duration := int(d.Seconds())

			if duration < 1 {
				duration = 1
			}

			icon, _ := args["icon"].(string)
			sound, ok := args["sound"].(string)

			if !ok {
				sound = "default"
			}

			channelTag, _ := args["channel"].(string)
			flowTag, _ := args["flow"].(string)

			notification := gotelemetry.NewNotification(title, message, icon, duration, sound)

			result := c.notificationProvider.SendNotification(notification, channelTag, flowTag)

			return newBooleanExpression(result, g.l, g.p), nil
		},
		map[string]callableArgument{
			"channel":  callableArgumentOptionalString,
			"flow":     callableArgumentOptionalString,
			"title":    callableArgumentString,
			"message":  callableArgumentString,
			"duration": callableArgumentOptionalString,
			"icon":     callableArgumentOptionalString,
			"sound":    callableArgumentOptionalString,
		},
		g.l,
		g.p,
	)
}

func (g *globalExpression) anomaly() expression {
	return newCallableExpression(
		"anomaly",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			data := args["data"].(numericArray)
			value := args["value"].(float64)

			mean := data.avg()
			stddev := data.stddev()

			return newBooleanExpression(math.Abs(value-mean) > 3*stddev, g.l, g.p), nil
		},
		map[string]callableArgument{
			"data":  callableArgumentNumericArray,
			"value": callableArgumentNumeric,
		},
		g.l,
		g.p,
	)
}

func mapFromYaml(from interface{}) interface{} {
	switch from.(type) {
	case map[interface{}]interface{}:
		result := map[string]interface{}{}

		for index, value := range from.(map[interface{}]interface{}) {
			result[index.(string)] = mapFromYaml(value)
		}

		return result

	case []interface{}:
		f := from.([]interface{})

		for index, value := range f {
			f[index] = mapFromYaml(value)
		}

		return f

	default:
		return from
	}
}

func (g *globalExpression) load() expression {
	return newCallableExpression(
		"load",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			format := strings.ToLower(args["format"].(string))
			path := args["path"].(string)

			data, err := ioutil.ReadFile(path)

			if err != nil {
				return nil, err
			}

			var result interface{}

			switch format {
			case "json":
				err = json.Unmarshal(data, &result)

			case "yaml":
				err = yaml.Unmarshal(data, &result)

			case "toml":
				_, err = toml.Decode(string(data), &result)

			default:
				return nil, errors.New(fmt.Sprintf("Unknown file format %s", format))
			}

			if err != nil {
				return nil, err
			}

			return expressionFromInterface(mapFromYaml(result), g.l, g.p)
		},
		map[string]callableArgument{
			"format": callableArgumentString,
			"path":   callableArgumentString,
		},
		g.l,
		g.p,
	)
}

func (g *globalExpression) spawn() expression {
	return newCallableExpression(
		"spawn",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			flowTag := args["tag"].(string)

			cfg := map[string]interface{}{
				"refresh":  int(args["refresh"].(float64)),
				"flow_tag": flowTag,
			}

			if script, ok := args["script"].(string); ok {
				cfg["script"] = script
			}

			if exec, ok := args["exec"].(string); ok {
				cfg["exec"] = exec
			}

			if ar, ok := args["args"].(map[string]interface{}); ok {
				cfg["args"] = ar
			}

			return nil, c.jobSpawner.SpawnJob(flowTag, "com.telemetryapp.process", cfg)
		},
		map[string]callableArgument{
			"exec":    callableArgumentOptionalString,
			"script":  callableArgumentOptionalString,
			"refresh": callableArgumentNumeric,
			"args":    callableArgumentOptionalMap,
			"tag":     callableArgumentString,
		},
		g.l,
		g.p,
	)
}

func execHttpRequest(method string, c *executionContext, args map[string]interface{}, l, p int) (expression, error) {
	reqUrl, err := url.Parse(args["url"].(string))
	if err != nil {
		return nil, err
	}

	query := reqUrl.Query()
	queryMap, _ := args["query"].(map[string]interface{})
	for k, v := range queryMap {
		query.Add(k, fmt.Sprintf("%v", v))
	}
	reqUrl.RawQuery = query.Encode()

	var contentType string

	var body io.Reader = nil
	paramsMap, _ := args["parameters"].(map[string]interface{})
	if len(paramsMap) > 0 {
		if isJson, ok := args["json"].(bool); ok && isJson {
			j, err := json.Marshal(paramsMap)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(j)
			contentType = "application/json"
		} else {
			params := url.Values{}
			for k, v := range paramsMap {
				params.Add(k, fmt.Sprintf("%v", v))
			}
			body = strings.NewReader(params.Encode())
			contentType = "application/x-www-form-urlencoded"
		}
	}

	req, err := http.NewRequest(method, reqUrl.String(), body)
	if err != nil {
		return nil, err
	}

	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}

	if auth, ok := args["auth"].(map[string]interface{}); ok && len(auth) > 0 {
		user, _ := auth["user"].(string)
		password, _ := auth["password"].(string)
		req.SetBasicAuth(user, password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respContentType := resp.Header.Get("Content-Type")

	var respBody interface{}
	if strings.HasPrefix(respContentType, "application/json") {
		err = json.Unmarshal(bodyBytes, &respBody)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(respContentType, "application/x-www-form-urlencoded") {
		respBody, err = url.ParseQuery(string(bodyBytes))
		if err != nil {
			return nil, err
		}
	} else {
		respBody = string(bodyBytes)
	}

	respHeader := map[string]interface{}{}
	for k, _ := range resp.Header {
		respHeader[k] = resp.Header.Get(k)
	}

	result := map[string]interface{}{
		"status_code": resp.StatusCode,
		"header":      respHeader,
		"body":        respBody,
	}

	return newMapExpression(result, l, p), nil
}

func (g *globalExpression) get() expression {
	return newCallableExpression(
		"get",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return execHttpRequest("GET", c, args, g.l, g.p)
		},
		map[string]callableArgument{
			"url":   callableArgumentString,
			"auth":  callableArgumentOptionalMap,
			"query": callableArgumentOptionalMap,
		},
		g.l,
		g.p,
	)
}

func (g *globalExpression) post() expression {
	return newCallableExpression(
		"post",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			return execHttpRequest("POST", c, args, g.l, g.p)
		},
		map[string]callableArgument{
			"url":        callableArgumentString,
			"auth":       callableArgumentOptionalMap,
			"query":      callableArgumentOptionalMap,
			"parameters": callableArgumentOptionalMap,
			"json":       callableArgumentOptionalBoolean,
		},
		g.l,
		g.p,
	)
}

func (g *globalExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := globalProperties[property]; ok {
		return f(g), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", g, property))
}

func (g *globalExpression) line() int {
	return g.l
}

func (g *globalExpression) position() int {
	return g.p
}

func (g *globalExpression) String() string {
	return fmt.Sprintf("Global()")
}
