package parser

import (
	"errors"
	"fmt"
)

type aslLexer struct {
	out      chan token
	current  token
	commands []command
	errors   []error
}

func newASLLexer(out chan token) *aslLexer {
	return &aslLexer{
		out:      out,
		commands: []command{},
		errors:   []error{},
	}
}

func (a *aslLexer) AddCommand(cmd command) {
	a.commands = append(a.commands, cmd)
}

func (a *aslLexer) Lex(lval *parserSymType) int {
	a.current = <-a.out

	lval.t = a.current

	return int(a.current.terminal)
}

func (a *aslLexer) Error(s string) {
	source := a.current.source

	if len(source) > 30 {
		source = source[0:29] + "..."
	}

	a.errors = append(a.errors, errors.New(fmt.Sprintf("Parse error: %s at line %d:%d (%q)", s, a.current.line, a.current.start, source)))
}
