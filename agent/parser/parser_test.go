package parser

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"sync"
	"testing"
)

type dummyNotificationProvider struct {
	notifications []gotelemetry.Notification
	channels      []string
}

func (d *dummyNotificationProvider) SendNotification(n gotelemetry.Notification, c string, f string) bool {
	d.notifications = append(d.notifications, n)
	d.channels = append(d.channels, c)

	return true
}

func newDummyNotificationProvider() *dummyNotificationProvider {
	return &dummyNotificationProvider{[]gotelemetry.Notification{}, []string{}}
}

var parserTestInitOnce = sync.Once{}

func testRunAndReturnErrors(s string) (map[string]interface{}, *dummyNotificationProvider, []error) {
	np := newDummyNotificationProvider()

	parserTestInitOnce.Do(func() {
		l := "/tmp/agent.sqlite3"
		aggregations.Init(&l, make(chan error, 99999))
	})

	commands, errs := Parse("test", s)

	if len(errs) > 0 {
		return nil, np, errs
	}

	if res, err := Run(np, nil, map[string]interface{}{"test": 10.0}, commands); err == nil {
		return res, np, nil
	} else {
		return res, np, []error{err}
	}
}

type testR map[string]interface{}
type testE []error

type parserTest struct {
	script    string
	condition interface{}
}

func runParserTests(tests map[string]parserTest, t *testing.T) {
	for index, test := range tests {
		res, np, errs := testRunAndReturnErrors(test.script)

		switch test.condition.(type) {
		case func(testR, testE) bool:

			if !test.condition.(func(testR, testE) bool)(res, errs) {
				for _, err := range errs {
					t.Errorf("Test %s -> error %s", index, err)
				}

				t.Errorf("Test %s fails condition: Got %#v", index, res)
			}

		case func(*dummyNotificationProvider) bool:

			if !test.condition.(func(*dummyNotificationProvider) bool)(np) {
				for _, err := range errs {
					t.Errorf("Test %s -> error %s", index, err)
				}

				t.Errorf("Test %s fails condition: Got %#v", index, res)
			}

		default:
			if len(errs) != 0 {
				for _, err := range errs {
					t.Errorf("Test %s -> error %s", index, err)
				}
			}

			if res["a"] != test.condition {
				t.Errorf("Unexpected result when running test %s: Wanted %T(%#v), got %T(%#v) instead", index, test.condition, test.condition, res["a"], res["a"])
			}
		}
	}
}

func TestBasicExpressions(t *testing.T) {
	tests := map[string]parserTest{
		"Comment":                  {`/* Test 1 2 3 * 3 */ a=123`, 123.0},
		"Numeric expression":       {"a=123", 123.0},
		"Addition":                 {"a=123+10", 133.0},
		"Multiplication":           {"a=10*10", 100.0},
		"Division":                 {"a=100/5", 20.0},
		"Subtraction":              {"a=132-10", 122.0},
		"Arithmetic precedence":    {"a=123+10*10", 223.0},
		"Parentheses":              {"a=(123+10)*10", 1330.0},
		"Unary Minus":              {"a=123+-10", 113.0},
		"Unary Minus + precedence": {"a=-(123+10)*10", -1330.0},
		"Variable assignment":      {"$a=10; a=$a+10", 20.0},
		"Addition with strings":    {`a="test"+10`, "test10"},
		"Addition with strings 2":  {`a=10+"test"`, func(res testR, errs testE) bool { return len(errs) == 1 }},
		"Variable transassignment": {`$a= series("cpu_load"); a= $a.last()`, func(res testR, errs testE) bool { _, ok := res["a"].(float64); return ok }},
	}

	runParserTests(tests, t)
}

func TestGlobalMethods(t *testing.T) {
	checkFloat := func(res testR, errs testE) bool {
		_, ok := res["a"].(float64)

		return ok
	}

	checkNotification := func(count int) func(*dummyNotificationProvider) bool {
		return func(np *dummyNotificationProvider) bool {
			return len(np.notifications) == 1
		}
	}

	tests := map[string]parserTest{
		"Global.now()":            {"a=now()", checkFloat},
		"Global.now() assignment": {"$a=now(); a=$a", checkFloat},
		"Global.notify()":         {`notify(channel:"123",title:"test",duration:"10s",message:"Hello")`, checkNotification(1)},
		"Global.arg()":            {`a=arg("test")`, 10.0},
		"Global.load()":           {`$a=load(format:"yaml",path:"parser_test_load.yaml"); a=$a.item("value")`, 123.0},
	}

	runParserTests(tests, t)
}

func TestSeries(t *testing.T) {
	checkFloat := func(res testR, errs testE) bool {
		_, ok := res["a"].(float64)

		return ok
	}

	checkArray := func(count int) func(res testR, errs testE) bool {
		return func(res testR, errs testE) bool {
			r, ok := res["a"].([]interface{})

			if ok {
				if len(r) != count {
					t.Errorf("Returned data should contain %d elements, but only %d found", count, len(r))
				}

				for _, rr := range r {
					if index, ok := rr.(float64); !ok {
						t.Errorf("Value %d in returned data is not a float.", index)
						return false
					}
				}
			}

			return ok
		}
	}

	tests := map[string]parserTest{
		"Series.last()":        {`a=series("cpu_load").last()+10`, checkFloat},
		"Series.aggregate()":   {`a=series("cpu_load").aggregate(func:"avg",interval:"10s",count:50)`, checkArray(50)},
		"Series.aggregate() 2": {`a=series("cpu_load").aggregate(func:"avg",interval:"10s",count:50).values.count()`, 50.0},
		"Series.aggregate() 3": {`a=series("cpu_load").aggregate(func:"avg",interval:"10s",count:10,end_time:"2010-01-02").values.count()`, 10.0},
		"Series.aggregate() 4": {`a=series("cpu_load").aggregate(func:"avg",interval:"10s",count:10,end_time:"2010-01-02 10:15:16").values.count()`, 10.0},
		"Series.avg()":         {`a=series("cpu_load").avg("10m")+10`, checkFloat},
		"Series.sum()":         {`a=series("cpu_load").avg("10m")+10`, checkFloat},
		"Series.count()":       {`a=series("cpu_load").count("10m")+10`, checkFloat},
		"Series.min()":         {`a=series("cpu_load").min("10m")+10`, checkFloat},
		"Series.max()":         {`a=series("cpu_load").max("10m")+10`, checkFloat},
		"Series.stddev()":      {`a=series("cpu_load").stddev("10m")+10`, checkFloat},
		"Series.trim(since)":   {`a=series("cpu_load").trim(since:"10m")`, nil},
		"Series.trim(count)":   {`a=series("cpu_load").trim(count:100)`, nil},
		"Series.items(count)":  {`a=series("nonexistent").items(100).values.count()`, 0.0},
	}

	runParserTests(tests, t)
}

func TestCounter(t *testing.T) {
	tests := map[string]parserTest{
		"Counter.set":       {`$a=counter("test123"); $a.set(100); a=counter("test123")`, 100.0},
		"Counter.increment": {`$a=counter("test123"); $a.set(100); $a.increment(100); a=counter("test123")`, 200.0},
		"Counter.reset":     {`$a=counter("test123"); $a.reset(); a=counter("test123")`, 0.0},
	}

	runParserTests(tests, t)
}

func TestBooleanAndLogicOperations(t *testing.T) {
	checkBool := func(expect bool) func(res testR, errs testE) bool {
		return func(res testR, errs testE) bool {
			r, ok := res["a"].(bool)

			if ok {
				return r == expect
			}

			return ok
		}
	}

	tests := map[string]parserTest{
		"Boolean false assignment": {"a=false", checkBool(false)},
		"Boolean true assignment":  {"a=true", checkBool(true)},
		"Boolean or 1":             {"a=true||true", checkBool(true)},
		"Boolean or 2":             {"a=true||false", checkBool(true)},
		"Boolean or 3":             {"a=false||true", checkBool(true)},
		"Boolean or 4":             {"a=false||false", checkBool(false)},
		"Boolean and 1":            {"a=true&&true", checkBool(true)},
		"Boolean and 2":            {"a=true&&false", checkBool(false)},
		"Boolean and 3":            {"a=false&&true", checkBool(false)},
		"Boolean and 4":            {"a=false&&false", checkBool(false)},
		"Equality 1":               {"a=true==true", checkBool(true)},
		"Equality 2":               {"a=true==false", checkBool(false)},
		"Equality 3":               {"a=false==true", checkBool(false)},
		"Equality 4":               {"a=false==false", checkBool(true)},
		"Equality 5":               {"a=10==10", checkBool(true)},
		"Equality 5.1":             {"a=10==11", checkBool(false)},
		"Equality 6":               {`a=10=="10"`, checkBool(true)},
		"Equality 6.1":             {`a=10=="11"`, checkBool(false)},
		"Equality 7":               {`a="test"=="test"`, checkBool(true)},
		"Equality 7.1":             {`a="test"=="test1"`, checkBool(false)},
		"Inequality 1":             {"a=true!=true", checkBool(false)},
		"Inequality 2":             {"a=true!=false", checkBool(true)},
		"Inequality 3":             {"a=false!=true", checkBool(true)},
		"Inequality 4":             {"a=false!=false", checkBool(false)},
		"Inequality 5":             {"a=10!=10", checkBool(false)},
		"Inequality 5.1":           {"a=10!=11", checkBool(true)},
		"Inequality 6":             {`a=10!="10"`, checkBool(false)},
		"Inequality 6.1":           {`a=10!="11"`, checkBool(true)},
		"Inequality 7":             {`a="test"!="test"`, checkBool(false)},
		"Inequality 7.1":           {`a="test"!="test1"`, checkBool(true)},
		"Greater than 1":           {`a=10>11`, checkBool(false)},
		"Greater than 2":           {`a=12>11`, checkBool(true)},
		"Greater than 3":           {`a=12>12`, checkBool(false)},
		"Greater than or equal 1":  {`a=10>=11`, checkBool(false)},
		"Greater than or equal 2":  {`a=12>=11`, checkBool(true)},
		"Greater than or equal 3":  {`a=12>=12`, checkBool(true)},
		"Less than 1":              {`a=10<11`, checkBool(true)},
		"Less than 2":              {`a=12<11`, checkBool(false)},
		"Less than 3":              {`a=12<12`, checkBool(false)},
		"Less than or equal 1":     {`a=10<=11`, checkBool(true)},
		"Less than or equal 2":     {`a=12<=11`, checkBool(false)},
		"Less than or equal 3":     {`a=12<=12`, checkBool(true)},
	}

	runParserTests(tests, t)
}

func TestIfThenElse(t *testing.T) {
	tests := map[string]parserTest{
		"If then":      {"if true==true{a=10}", 10.0},
		"If then else": {"if false==true{a=10}else{a=20}", 20.0},
	}

	runParserTests(tests, t)
}

func TestArrays(t *testing.T) {
	checkArray := func(expect []interface{}) func(res testR, errs testE) bool {
		return func(res testR, errs testE) bool {
			r := res["a"].([]interface{})

			if len(r) != len(expect) {
				return false
			}

			for index, value := range expect {
				if r[index] != value {
					return false
				}
			}

			return true
		}
	}

	tests := map[string]parserTest{
		"Immediate array": {"a = [10, 20, 30]", checkArray([]interface{}{10.0, 20.0, 30.0})},
		"Array element":   {"$a = [10, 20, 30]; a = $a.item(1)", 20.0},
		"Array sum":       {"$a = [10, 20, 30]; a = $a.sum()", 60.0},
		"Array min":       {"$a = [10, 20, 30]; a = $a.min()", 10.0},
		"Array max":       {"$a = [10, 20, 30]; a = $a.max()", 30.0},
		"Array avg":       {"$a = [10, 20, 30]; a = $a.avg()", 20.0},
		"Array count":     {"$a = [10, 20, 30]; a = $a.count()", 3.0},
		"Array stddev":    {"$a = [10, 20, 30]; a = $a.stddev()", 8.16496580927726},
	}

	runParserTests(tests, t)
}

func TestMaps(t *testing.T) {
	compareArray := func(expect []interface{}, r []interface{}) bool {
		if len(r) != len(expect) {
			return false
		}

		for index, value := range expect {
			if r[index] != value {
				return false
			}
		}

		return true
	}

	compareMaps := func(expect map[string]interface{}) func(res testR, errs testE) bool {
		return func(res testR, errs testE) bool {
			r := res["a"].(map[string]interface{})

			for i, v := range expect {
				if vv, ok := r[i]; ok {
					if vvv, ok := v.([]interface{}); ok {
						if !compareArray(vvv, vv.([]interface{})) {
							return false
						} else {

						}
					} else {
						if v != vv {
							return false
						}
					}
				} else {
					return false
				}
			}

			for i, v := range r {
				if vv, ok := expect[i]; ok {
					if vvv, ok := vv.([]interface{}); ok {
						if !compareArray(vvv, v.([]interface{})) {
							return false
						} else {

						}
					} else {
						if v != vv {
							return false
						}
					}
				} else {
					return false
				}
			}

			return true
		}
	}

	tests := map[string]parserTest{
		"Immediate map":  {`a = {a:10, b:"test", c:[10, 20, 30]}`, compareMaps(map[string]interface{}{"a": 10.0, "b": "test", "c": []interface{}{10.0, 20.0, 30.0}})},
		"Map assignment": {`$a = {a:10, b:"test", c:[10, 20, 30]}; a = $a`, compareMaps(map[string]interface{}{"a": 10.0, "b": "test", "c": []interface{}{10.0, 20.0, 30.0}})},
		"Map item":       {`$a = {a:10, b:"test", c:[10, 20, 30]}; a = $a.item("a")`, 10.0},
		"Map count":      {`$a = {a:10, b:"test", c:[10, 20, 30]}; a = $a.count()`, 3.0},
		"Map set":        {`$a = {a:10, b:"test", c:[10, 20, 30]}; $a.set(index:"a", value:10+10); a = $a.item("a")`, 20.0},
		"Map and array":  {`$a = [ { color: "red", values: series("sampleseries").sum("10s") } ]; a = $a.item(0)`, compareMaps(map[string]interface{}{"values": 0.0, "color": "red"})},
	}

	runParserTests(tests, t)
}

func TestAnomaly(t *testing.T) {
	tests := map[string]parserTest{
		"Anomaly (true)":  {"$a = [10, 20, 30]; a = anomaly(data:$a, value:10000)", true},
		"Anomaly (false)": {"$a = [10, 20, 30]; a = anomaly(data:$a, value:10)", false},
	}

	runParserTests(tests, t)
}

func TestWhile(t *testing.T) {
	tests := map[string]parserTest{
		"While loop": {"$a = 10; while $a < 20 { $a = $a + 1 } a = $a", 20.0},
	}

	runParserTests(tests, t)
}

func TestExcel(t *testing.T) {
	tests := map[string]parserTest{
		"Excel 1": {`a = excel("excel_test.xlsx").cells(ranges:"A2").item(0)`, 10.0},
		"Excel 2": {`a = excel("excel_test.xlsx").cells(ranges:"C10:E10").sum()`, 90.0},
	}

	runParserTests(tests, t)
}

func TestGet(t *testing.T) {
	tests := map[string]parserTest{
		"Get 1": {`a = get(url:"http://jsonplaceholder.typicode.com/users")`, func(res testR, errs testE) bool {
			a := res["a"].(map[string]interface{})
			return (a["status_code"] == 200.0 &&
				len(a["body"].([]interface{})) == 10)
		}},
		"Get 2": {`a = get(url:"http://jsonplaceholder.typicode.com/users", query:{id:2})`, func(res testR, errs testE) bool {
			a := res["a"].(map[string]interface{})
			return (a["status_code"] == 200.0 &&
				len(a["body"].([]interface{})) == 1)
		}},
	}

	runParserTests(tests, t)
}

func TestPost(t *testing.T) {
	tests := map[string]parserTest{
		"Post 1": {`a = post(url:"http://jsonplaceholder.typicode.com/posts", parameters:{title:"blah",body:"foobar",userId:1}, json:true)`, func(res testR, errs testE) bool {
			a := res["a"].(map[string]interface{})
			body := a["body"].(map[string]interface{})
			return (a["status_code"] == 200.0 &&
				body["title"].(string) == "blah" &&
				body["body"].(string) == "foobar")
		}},
		"Post 2": {`a = post(url:"http://jsonplaceholder.typicode.com/posts", parameters:{title:"blah",body:"foobar",userId:1})`, func(res testR, errs testE) bool {
			a := res["a"].(map[string]interface{})
			body := a["body"].(map[string]interface{})
			return (a["status_code"] == 200.0 &&
				body["title"].(string) == "blah" &&
				body["body"].(string) == "foobar")
		}},
	}

	runParserTests(tests, t)
}
