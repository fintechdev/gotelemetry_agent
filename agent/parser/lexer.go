package parser

import (
	"unicode/utf8"
)

const (
	symbol      = "+-/*:.,()[]"
	numeric     = "0123456789."
	whitespace  = " \n\t"
	openParens  = "("
	closeParens = ")"
	stringStart = `"`
	stringEnd   = `\"`

	identifier = symbol + whitespace + openParens + closeParens + stringStart + stringEnd
)

var symbolMap = map[rune]terminal{
	'+': T_PLUS,
	'-': T_MINUS,
	'/': T_DIVIDE,
	'*': T_MULTIPLY,
	':': T_COLON,
	'.': T_DOT,
	',': T_COMMA,
	'(': T_OPEN_PARENS,
	')': T_CLOSE_PARENS,
	'[': T_OPEN_BRACKET,
	']': T_CLOSE_BRACKET,
}

func lexASL(name string, source string) chan token {
	l := lex(name, source, aslInitial)

	return l.out
}

func aslInitial(l *lexer) stateFn {
	// Each all whitespace

	l.acceptRun(whitespace, false)
	l.ignore()

	// Check for eof

	if l.peek() == eofRune {
		return nil
	}

	// Check for symbols

	if r := l.accept(symbol, false); r != utf8.RuneError {
		if t, ok := symbolMap[r]; ok {
			l.emit(t)
			return aslInitial
		} else {
			l.errorf("Unknown symbol %q", r)
			return nil
		}
	}

	// Check for numbers

	if l.accept(numeric, false) != utf8.RuneError {
		l.acceptRun(numeric, false)
		l.emit(T_NUMBER)

		return aslInitial(l)
	}

	// Check for strings

	if l.accept(stringStart, false) != utf8.RuneError {
		l.ignore() // We don't want the string delimiter, only the string
		return aslString
	}

	// Must be an identifier

	l.acceptRun(identifier, true)
	l.emit(T_IDENTIFIER)

	return aslInitial
}

func aslString(l *lexer) stateFn {
	l.acceptRun(stringEnd, true)

	r := l.next()

	switch r {
	case '\\':
		l.next() // Eat the next run, since we're escaping it

		return aslString

	case '"':
		l.backup() // Go back one; we don't want to include the quotation mark
		l.emit(T_STRING)

		l.next()
		l.ignore()

		return aslInitial

	default:
		l.errorf("Unterminated string")
		return nil
	}

}
