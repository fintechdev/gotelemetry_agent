package parser

import (
	"errors"
	"fmt"
	"testing"
)

func aslTestingParse(s string) ([]token, error) {
	out := lexASL("Lex tester", s)

	tokes := []token{}

	for {
		toke := <-out

		switch toke.terminal {
		case eofTerminal:
			return tokes, nil

		case errorTerminal:
			return nil, errors.New(fmt.Sprintf("Error [%d:%d] %s", toke.line, toke.start, toke.source))

		default:
			tokes = append(tokes, toke)
		}
	}
}

func TestASLLexer(t *testing.T) {
	s := `value:series("test").aggregate("avg","7d","10m")`

	tokes, err := aslTestingParse(s)

	if err != nil {
		t.Errorf("Unable to parse valid expression: %s", err)
	}

	if !compareTokens(tokes, []token{token{terminal: 57348, line: 1, start: 0, source: "value"}, token{terminal: 57359, line: 1, start: 5, source: ":"}, token{terminal: 57348, line: 1, start: 6, source: "series"}, token{terminal: 57360, line: 1, start: 12, source: "("}, token{terminal: 57346, line: 1, start: 14, source: "test"}, token{terminal: 57361, line: 1, start: 19, source: ")"}, token{terminal: 57357, line: 1, start: 20, source: "."}, token{terminal: 57348, line: 1, start: 21, source: "aggregate"}, token{terminal: 57360, line: 1, start: 30, source: "("}, token{terminal: 57346, line: 1, start: 32, source: "avg"}, token{terminal: 57356, line: 1, start: 36, source: ","}, token{terminal: 57346, line: 1, start: 38, source: "7d"}, token{terminal: 57356, line: 1, start: 41, source: ","}, token{terminal: 57346, line: 1, start: 43, source: "10m"}, token{terminal: 57361, line: 1, start: 47, source: ")"}, token{terminal: 57377, line: 1, start: 48, source: ""}}) {
		t.Errorf("Unexpected lexer output: %#v", tokes)
	}
}
