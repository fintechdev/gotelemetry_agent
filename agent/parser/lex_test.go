package parser

import (
	"errors"
	"fmt"
	"testing"
	"unicode/utf8"
)

// A simple lexer that takes the format
// NUMBER [+-/*] NUMBER

const (
	testNumberTerminal terminal = iota
	testOperatorTerminal
)

func lexTestOperatorState(l *lexer) stateFn {
	l.acceptRun(" \t\n", false)
	l.ignore()

	r := l.accept("+-/*", false)

	if r == utf8.RuneError {
		l.errorf("Expected one of `+`, `-`, `*`, `/`, found %q instead", l.next())
		return nil
	}

	if r == eofRune {
		return nil
	}

	l.emit(testOperatorTerminal)

	return lexTestInitialState
}

func lexTestInitialState(l *lexer) stateFn {
	l.acceptRun(" \t\n", false)
	l.ignore()

	r := l.accept("+-0123456789.", false)

	if r == utf8.RuneError {
		l.errorf("Expected number, found %q instead", l.next())
		return nil
	}

	if r == eofRune {
		return nil
	}

	l.acceptRun("0123456789.", false)

	if l.length > 0 {
		l.emit(testNumberTerminal)
		return lexTestOperatorState
	}

	return nil
}

func lexTestingParse(s string) ([]token, error) {
	l := lex("Lex tester", s, lexTestInitialState)

	tokes := []token{}

	for {
		toke := <-l.out

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

func compareTokens(left, right []token) bool {
	if len(left) != len(right) {
		return false
	}

	for index, l := range left {
		if l != right[index] {
			return false
		}
	}

	return true
}

func TestLexer(t *testing.T) {
	tokes, err := lexTestingParse("123 + 192 / -21")

	if err != nil {
		t.Error(err)
	}

	if !compareTokens(tokes, []token{token{terminal: 0, line: 1, start: 0, source: "123"}, token{terminal: 1, line: 1, start: 4, source: "+"}, token{terminal: 0, line: 1, start: 6, source: "192"}, token{terminal: 1, line: 1, start: 10, source: "/"}, token{terminal: 0, line: 1, start: 12, source: "-21"}, token{terminal: 57376, line: 1, start: 15, source: ""}}) {
		t.Errorf("Unexpected lexer output: %#v", tokes)
	}

	tokes, err = lexTestingParse("123 ++192 / -21")

	if err != nil {
		t.Error(err)
	}

	if !compareTokens(tokes, []token{token{terminal: 0, line: 1, start: 0, source: "123"}, token{terminal: 1, line: 1, start: 4, source: "+"}, token{terminal: 0, line: 1, start: 5, source: "+192"}, token{terminal: 1, line: 1, start: 10, source: "/"}, token{terminal: 0, line: 1, start: 12, source: "-21"}, token{terminal: 57376, line: 1, start: 15, source: ""}}) {
		t.Errorf("Unexpected lexer output: %#v", tokes)
	}

	tokes, err = lexTestingParse("SAMURAI 123")

	if err == nil {
		t.Errorf("Lexer accepted an invalid string")
	}
}
