package parser

import (
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"testing"
	"time"
)

func testRun(s string, t *testing.T) map[string]interface{} {
	commands, errs := Parse("test", s)

	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}

		return nil
	}

	if res, err := Run(commands); err == nil {
		return res
	} else {
		t.Error(err)

		return nil
	}
}

func testRunAndReturnErrors(s string) []error {
	commands, errs := Parse("test", s)

	if len(errs) > 0 {
		return errs
	}

	if _, err := Run(commands); err == nil {
		return nil
	} else {
		return []error{err}
	}
}

func TestNumericExpression(t *testing.T) {
	l := "/tmp/telemetry.sqlite"
	aggregations.Init(&l, make(chan error, 99999))

	res := testRun("a:123", t)

	if res["a"] != 123.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}
}

func TestArithmeticExpressions(t *testing.T) {
	l := "/tmp/telemetry.sqlite"
	aggregations.Init(&l, make(chan error, 99999))

	res := testRun("a:123+10", t)

	if res["a"] != 133.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}

	res = testRun("a:10*10", t)

	if res["a"] != 100.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}

	res = testRun("a:100/5", t)

	if res["a"] != 20.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}

	res = testRun("a:123-10", t)

	if res["a"] != 113.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}

	res = testRun("a:123+10*10", t)

	if res["a"] != 223.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}

	res = testRun("a:(123+10)*10", t)

	if res["a"] != 1330.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}
}

func TestUnaryExpressions(t *testing.T) {
	l := "/tmp/telemetry.sqlite"
	aggregations.Init(&l, make(chan error, 99999))

	res := testRun("a:123+-10", t)

	if res["a"] != 113.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}

	res = testRun("a:-(123+10)*10", t)

	if res["a"] != -1330.0 {
		t.Errorf("Unexpected expression result: %v", res)
	}
}

func TestArithmeticDeviance(t *testing.T) {
	l := "/tmp/telemetry.sqlite"
	aggregations.Init(&l, make(chan error, 99999))

	err := testRunAndReturnErrors(`a:"test"+10`)

	if err == nil {
		t.Error("Numeric operations can be performed with non-numeric values")
	}
}

func TestVariableAssignment(t *testing.T) {
	l := "/tmp/telemetry.sqlite"
	aggregations.Init(&l, make(chan error, 99999))

	res := testRun("$a: 10 a:$a+10", t)

	if res["a"] == 20 {
		t.Errorf("Unexpected expression result: %v", res)
	}
}

func TestGlobalMethods(t *testing.T) {
	l := "/tmp/telemetry.sqlite"
	aggregations.Init(&l, make(chan error, 99999))

	now := time.Now().Unix()

	res := testRun("a: now()", t)

	if n, ok := res["a"].(float64); !ok || int64(n) < now {
		t.Errorf("Unexpected expression result: %v", res)
	}
}
