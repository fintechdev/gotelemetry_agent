package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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
	paramsMap, _ := args["body"].(map[string]interface{})
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
			"url":   callableArgumentString,
			"auth":  callableArgumentOptionalMap,
			"query": callableArgumentOptionalMap,
			"body":  callableArgumentOptionalMap,
			"json":  callableArgumentOptionalBoolean,
		},
		g.l,
		g.p,
	)
}
