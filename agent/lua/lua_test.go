package lua

import (
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"testing"
	"time"
)

type expectsError bool
type resultValidator func(*testing.T, error, map[string]interface{}) bool

var shouldError = expectsError(true)
var shouldNotError = expectsError(false)

type test struct {
	name   string
	source string
	result interface{}
}

func compareValue(left, right interface{}) bool {
	switch left.(type) {
	case map[string]interface{}:
		if r, ok := right.(map[string]interface{}); ok {
			return compareMap(left.(map[string]interface{}), r)
		}

		return false

	case []interface{}:
		if r, ok := right.([]interface{}); ok {
			return compareArray(left.([]interface{}), r)
		}

		return false

	default:
		return left == right
	}
}

func compareMap(left, right map[string]interface{}) bool {
	for k, v := range left {
		if vv, ok := right[k]; ok {
			if !compareValue(v, vv) {
				return false
			}

			continue
		}

		return false
	}

	return true
}

func compareArray(left, right []interface{}) bool {
	if len(left) != len(right) {
		return false
	}

	for i, v := range left {
		if !compareValue(v, right[i]) {
			return false
		}
	}

	return true
}

type dummyNotificationProvider struct{}

func (d *dummyNotificationProvider) SendNotification(n gotelemetry.Notification, c string, f string) bool {
	return true
}

func runTests(t *testing.T, tests []test) {
	for _, tt := range tests {
		output, err := Exec(tt.source, &dummyNotificationProvider{}, map[string]interface{}{"test": 123})

		switch tt.result.(type) {
		case expectsError:
			v := tt.result.(expectsError)

			if v && err == nil {
				t.Errorf("Test %s should return an error, but does not.", tt.name)
			}

			if !v && err != nil {
				t.Errorf("Test %s should not return an error, but returned `%s`.", tt.name, err)
			}

		case resultValidator:
			if !tt.result.(resultValidator)(t, err, output) {
				t.Errorf("Test %s fails result validation.", tt.name)
			}

		default:
			if err != nil {
				t.Errorf("Test %s should not return an error, but returned `%s`.", tt.name, err)
			}

			if !compareValue(tt.result, output) {
				t.Errorf("Test %s fails result validation. Expected `%#v`, but got `%#v` instead", tt.name, tt.result, output)
			}
		}
	}
}

/*
func TestMongo(t *testing.T) {
	runTests(
		t,
		[]test{
			{"Connect to a Mongo database", `local mongo = require("telemetry/mongodb"); local session = mongo.open("mongodb://localhost:27017/local"); session.close();`, map[string]interface{}{}},
			{"Retrieve list of live servers", `local mongo = require("telemetry/mongodb"); local session = mongo.open("mongodb://localhost:27017/local"); output.data = session.live_servers(); session.close();`, map[string]interface{}{"data": []interface{}{"localhost:27017"}}},
			{"Get a DB", `local mongo = require("telemetry/mongodb"); local session = mongo.open("mongodb://localhost:27017/local"); local db = session.db("local"); output.data = db.collections(); session.close();`, map[string]interface{}{"data": []interface{}{"startup_log", "system.indexes"}}},
			{"Get a Collection name", `local mongo = require("telemetry/mongodb"); local session = mongo.open("mongodb://localhost:27017/local"); local db = session.db("local"); local collection = db.collection("startup_log"); output.data = collection.name(); session.close();`, map[string]interface{}{"data": "startup_log"}},
			{"Get a Collection and perform a search", `local mongo = require("telemetry/mongodb"); local session = mongo.open("mongodb://localhost:27017/local"); local db = session.db("local"); local collection = db.collection("startup_log"); local query = {}; query["cmdLine.net.bindIp"] = "127.0.0.1"; output.data = #(collection.query(query, 0, 1)); session.close();`, map[string]interface{}{"data": 1.0}},
			{"Perform a search and get a count", `local mongo = require("telemetry/mongodb"); local session = mongo.open("mongodb://localhost:27017/local"); local db = session.db("local"); local collection = db.collection("startup_log"); local query = {}; query["cmdLine.net.bindIp"] = "127.0.0.1"; output.data = collection.count(query, 0, 10); session.close();`, map[string]interface{}{"data": 10.0}},
			{"Run a command", `local mongo = require("telemetry/mongodb"); local session = mongo.open("mongodb://localhost:27017/local"); local db = session.db("admin"); output.data = db.command({ ping = 1 }); session.close();`, map[string]interface{}{"data": map[string]interface{}{"ok": 1.0}}},
		},
	)
}
*/
func TestRunScript(t *testing.T) {
	runTests(
		t,
		[]test{
			{"Run a script", `output.test = "123"`, map[string]interface{}{"test": "123"}},
			{"Attempt to overwrite a global", `output = "123"`, shouldError},
			{"Access args", `output.out = args.test`, map[string]interface{}{"out": 123.0}},
			{"Output an array", `output.out = { 1 , 2 , 3 }`, map[string]interface{}{"out": map[string]interface{}{"1": 1.0, "2": 2.0, "3": 3.0}}},
		},
	)
}

func TestJSON(t *testing.T) {
	runTests(
		t,
		[]test{
			{"JSON Encode String", `local json = require("telemetry/json"); output.out = json.encode("Test 1 2 3")`, map[string]interface{}{"out": `"Test 1 2 3"`}},
			{"JSON Encode Table", `local json = require("telemetry/json"); output.out = json.encode({a = 123})`, map[string]interface{}{"out": `{"a":123}`}},
			{"JSON Decode", `local json = require("telemetry/json"); output.out = json.decode("{\"a\":123}")`, map[string]interface{}{"out": map[string]interface{}{"a": 123.0}}},
		},
	)
}

func TestXML(t *testing.T) {
	runTests(
		t,
		[]test{
			{"XML Encode String", `local xml = require("telemetry/xml"); output.out = xml.encode("Test 1 2 3")`, shouldError},
			{"XML Encode Table", `local xml = require("telemetry/xml"); output.out = xml.encode({a = 123})`, map[string]interface{}{"out": "<a>123</a>"}},
			{"XML Decode", `local xml = require("telemetry/xml"); output.out = xml.decode("<note type=\"123\"><to>Tove</to><from>Jani</from><heading>Reminder</heading><body>Don't forget me this weekend!</body></note>")`, map[string]interface{}{"out": map[string]interface{}{"note": map[string]interface{}{"body": "Don't forget me this weekend!", "to": "Tove", "from": "Jani", "heading": "Reminder"}}}},
		},
	)
}

func TestHTTP(t *testing.T) {
	runTests(
		t,
		[]test{
			{"HTTP GET", `local http = require("telemetry/http"); output.out = http.get("https://raw.githubusercontent.com/telemetryapp/gotelemetry_agent/6efa4be88b4072a72f4a0d47cb15ca4e15263663/VERSION")`, map[string]interface{}{"out": "2.2.4"}},
			{"HTTP POST", `local http = require("telemetry/http"); output.out = http.post("http://jsonplaceholder.typicode.com/posts", "{\"title\":\"blah\",\"body\":\"foobar\",\"userId\":1}", {["Content-Type"] = "application/json"})`, map[string]interface{}{"out": "{\n  \"title\": \"blah\",\n  \"body\": \"foobar\",\n  \"userId\": 1,\n  \"id\": 101\n}"}},
			{"HTTP Authenticated GET", `local http = require("telemetry/http"); output.out = http.get("https://httpbin.org/basic-auth/atestusr/secretpass", "atestusr", "secretpass")`, map[string]interface{}{"out": "{\n  \"authenticated\": true, \n  \"user\": \"atestusr\"\n}\n"}},
			{"HTTP Authenticated Custom GET", `local http = require("telemetry/http"); output.out = http.custom("GET", "https://httpbin.org/basic-auth/atestusr/secretpass", "", "atestusr", "secretpass")`, map[string]interface{}{"out": "{\n  \"authenticated\": true, \n  \"user\": \"atestusr\"\n}\n"}},
			{"HTTP Custom POST", `local http = require("telemetry/http"); output.out = http.custom("POST", "http://jsonplaceholder.typicode.com/posts", "{\"title\":\"blah\",\"body\":\"foobar\",\"userId\":1}", {["Content-Type"] = "application/json"})`, map[string]interface{}{"out": "{\n  \"title\": \"blah\",\n  \"body\": \"foobar\",\n  \"userId\": 1,\n  \"id\": 101\n}"}},
		},
	)
}

func TestRegex(t *testing.T) {
	script := `
	local regex = require("goluago/regexp")

	function urlencode(str)
		if (str) then

			local res = ""

			for i = 1, string.len(str) do
				local ch = string.sub(str, i, i)

				if regex.match("[^A-Za-z0-9\\_.~]", ch) then
					res = res .. string.format("%%%02X", string.byte(ch))
				else
					res = res .. ch
				end
			end

			str = res
		end

		return str
	end

	output.test = urlencode("Test+Meé")
	`

	runTests(
		t,
		[]test{
			{"Regex replace", script, map[string]interface{}{"test": "Test%2BMe%C3%A9"}},
		},
	)
}

func TestSeries(t *testing.T) {
	l := "/tmp/agent.bolt"
	ttl := "1h"
	aggregations.Init(nil, &l, &ttl, make(chan error, 99999))

	ts := float64(time.Now().Unix() + 10)
	tss := fmt.Sprintf("%g", ts)

	runTests(
		t,
		[]test{
			{"Series", `local st = require("telemetry/storage"); st.series("test")`, shouldNotError},
			{"Series name", `local st = require("telemetry/storage"); output.out = st.series("test").name()`, map[string]interface{}{"out": "test"}},
			{"Series trim since by timestamp", `local st = require("telemetry/storage"); st.series("test").trimSince(os.time() - (60 * 2))`, shouldNotError},
			{"Series trim since by duration", `local st = require("telemetry/storage"); st.series("test").trimSince("2m")`, shouldNotError},
			{"Series trim by count", `local st = require("telemetry/storage"); st.series("test").trimCount(30)`, shouldNotError},
			{"Series push", `local st = require("telemetry/storage"); st.series("test").push(123)`, shouldNotError},
			{"Series pop", `local st = require("telemetry/storage"); s = st.series("test"); s.push(125, "` + tss + `"); output.out = s.pop(true)`, map[string]interface{}{"out": map[string]interface{}{"value": 125.0, "ts": ts}}},
			{"Series last", `local st = require("telemetry/storage"); s = st.series("test"); s.push(126, "` + tss + `"); output.out = s.last()`, map[string]interface{}{"out": map[string]interface{}{"value": 126.0, "ts": ts}}},
			{"Series compute by timestamp", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").compute(st.Functions.SUM, os.time() - 10000000, os.time() + 10000)`, shouldNotError},
			{"Series compute by interval", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").compute(st.Functions.SUM, "6m")`, shouldNotError},
			{"Series get raw items", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").items(10)`, shouldNotError},
			{"Series aggregate", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").aggregate(st.Functions.SUM, 60, 10, os.time())`, resultValidator(func(t *testing.T, err error, res map[string]interface{}) bool {
				return err == nil && len(res["out"].([]interface{})) == 10
			})},
			{"Series aggregate no end time", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").aggregate(st.Functions.SUM, 60, 10)`, resultValidator(func(t *testing.T, err error, res map[string]interface{}) bool {
				return err == nil && len(res["out"].([]interface{})) == 10
			})},
			{"Series aggregate string interval", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").aggregate(st.Functions.SUM, "1m", 10)`, resultValidator(func(t *testing.T, err error, res map[string]interface{}) bool {
				return err == nil && len(res["out"].([]interface{})) == 10
			})},
			{"Series values", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").aggregate(st.Functions.SUM, 60, 10, os.time()).values()`, resultValidator(func(t *testing.T, err error, res map[string]interface{}) bool {
				return err == nil && len(res["out"].([]interface{})) == 10
			})},
			{"Series timestamps", `local os = require("os"); local st = require("telemetry/storage"); output.out = st.series("test").aggregate(st.Functions.SUM, 60, 10, os.time()).ts()`, resultValidator(func(t *testing.T, err error, res map[string]interface{}) bool {
				return err == nil && len(res["out"].([]interface{})) == 10
			})},
		},
	)
}

func TestCounter(t *testing.T) {
	runTests(
		t,
		[]test{
			{"Counter", `local st = require("telemetry/storage"); st.counter("test")`, shouldNotError},
			{"Counter Value", `local st = require("telemetry/storage"); output.out = st.counter("test").value()`, shouldNotError},
			{"Counter Set", `local st = require("telemetry/storage"); c = st.counter("test"); c.set(10); output.out = st.counter("test").value()`, map[string]interface{}{"out": 10.0}},
			{"Counter Increment", `local st = require("telemetry/storage"); c = st.counter("test"); c.set(10); c.increment(1); output.out = st.counter("test").value()`, map[string]interface{}{"out": 11.0}},
		},
	)
}

func TestNotifications(t *testing.T) {
	runTests(
		t,
		[]test{
			{"Notifications", `local n = require("telemetry/notifications"); output.out = n.post("Channel", "Tag", "Title", "Message", 10)`, map[string]interface{}{"out": true}},
		},
	)
}

func TestExcel(t *testing.T) {
	runTests(
		t,
		[]test{
			{"Excel", `local excel = require("telemetry/excel"); output.out = tonumber(excel.import("excel_test.xlsx")[1][4][1])`, map[string]interface{}{"out": 10.0}},
		},
	)
}

func TestErrors(t *testing.T) {

	source := `

	local a;

	invalid code

	`

	source2 := `

	local st = require("telemetry/storage")
	st.series(11)

	`

	runTests(
		t,
		[]test{
			{"Syntax errors", source, shouldError},
			{"Runtime errors", source2, shouldError},
		},
	)
}
