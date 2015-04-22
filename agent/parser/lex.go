package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type stateFn func(l *lexer) stateFn

const eofRune = -1

type terminal int

const (
	eofTerminal   terminal = -1
	errorTerminal terminal = -2
)

type token struct {
	terminal terminal
	line     int
	start    int
	source   string
}

type lexer struct {
	name    string     // Name of the lexer
	source  string     // Source code
	start   int        // Start of the current token
	pos     int        // Current position
	length  int        // Length of the rune at the current position
	lastPos int        // Position of the previous token
	out     chan token // Output channel
}

func lex(name, source string, initial stateFn) *lexer {
	l := &lexer{
		name:   name,
		source: source,
		out:    make(chan token, 2),
	}

	go l.run(initial)

	return l
}

func (l *lexer) next() rune {
	// Make sure we're not already at the end of the input

	if l.pos >= len(l.source) {
		return eofRune
	}

	// Find the next rune
	r, length := utf8.DecodeRuneInString(l.source[l.pos:])

	// Move to the beginning of the next rune
	l.pos += length
	l.length = length

	// Return it
	return r
}

func (l *lexer) backup() {
	// Go back to the last rune
	l.pos -= l.length
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

func (l *lexer) peekAccept(valid string, negate bool) rune {
	r := l.accept(valid, negate)

	if r != utf8.RuneError {
		l.backup()
	}

	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string, negate bool) rune {
	r := l.next()

	if r == eofRune || ((strings.IndexRune(valid, r) >= 0) != negate) {
		return r
	}

	l.backup()
	return utf8.RuneError
}

func (l *lexer) acceptRun(valid string, negate bool) {
	for r := l.accept(valid, negate); r != eofRune && r != utf8.RuneError; {
		r = l.accept(valid, negate)
	}
}

func (l *lexer) lineNumber() int {
	return strings.Count(l.source[:l.pos], "\n") + 1
}

func (l *lexer) errorf(format string, args ...interface{}) {
	l.out <- token{terminal: errorTerminal, line: l.lineNumber(), start: l.start, source: fmt.Sprintf(format, args...)}
	l.start = l.pos
}

func (l *lexer) current() string {
	return l.source[l.start:l.pos]
}

func (l *lexer) emit(v terminal) {
	if l.pos > l.start {
		l.out <- token{terminal: v, line: l.lineNumber(), start: l.start, source: l.source[l.start:l.pos]}
		l.start = l.pos
	}
}

func (l *lexer) run(initial stateFn) {
	for state := initial; state != nil; {
		state = state(l)
	}

	l.out <- token{terminal: eofTerminal, line: l.lineNumber(), start: l.start, source: ""}
	close(l.out)
}
