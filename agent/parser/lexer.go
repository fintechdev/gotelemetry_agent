package parser

import (
	"unicode/utf8"
)

const (
	symbol        = "+-/*:.,()[]=&|!{}\n;<>"
	numeric       = "0123456789."
	whitespace    = " \t"
	openParens    = "("
	closeParens   = ")"
	stringStart   = `"`
	stringEnd     = `\"`
	variableStart = "$"

	identifier = symbol + whitespace + openParens + closeParens + stringStart + stringEnd + variableStart
)

var symbolMap = map[string]terminal{
	"+":  T_PLUS,
	"-":  T_MINUS,
	"/":  T_DIVIDE,
	"*":  T_MULTIPLY,
	":":  T_COLON,
	"=":  T_ASSIGN,
	".":  T_DOT,
	",":  T_COMMA,
	"(":  T_OPEN_PARENS,
	")":  T_CLOSE_PARENS,
	"[":  T_OPEN_BRACKET,
	"]":  T_CLOSE_BRACKET,
	"==": T_EQUAL,
	"&&": T_AND,
	"||": T_OR,
	"!":  T_NEGATE,
	"!=": T_NOT_EQUAL,
	"{":  T_OPEN_BRACE,
	"}":  T_CLOSE_BRACE,
	"\n": T_TERMINATOR,
	";":  T_TERMINATOR,
	">":  T_GREATER_THAN,
	"<":  T_LESS_THAN,
	">=": T_GREATER_THAN_OR_EQUAL,
	"<=": T_LESS_THAN_OR_EQUAL,
	"/*": T_COMMENT,
}

var identifierMap = map[string]terminal{
	"false": T_FALSE,
	"true":  T_TRUE,
	"if":    T_IF,
	"else":  T_ELSE,
	"while": T_WHILE,
	"null":  T_NULL,
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
		return aslSymbol
	}

	// Check for numbers

	if l.accept(numeric, false) != utf8.RuneError {
		l.acceptRun(numeric, false)
		l.emit(T_NUMBER)

		return aslInitial
	}

	// Check for strings

	if l.accept(stringStart, false) != utf8.RuneError {
		l.ignore() // We don't want the string delimiter, only the string
		return aslString
	}

	// Check for variable names

	if l.accept(variableStart, false) != utf8.RuneError {
		return aslVariableName
	}

	// Must be an identifier

	l.acceptRun(identifier, true)

	if t, ok := identifierMap[l.current()]; ok {
		l.emit(t)
	} else {
		l.emit(T_IDENTIFIER)
	}

	return aslInitial
}

func aslComment(l *lexer) stateFn {
	for {
		l.acceptRun("*", true)
		l.next()

		switch l.next() {
		case '/':
			l.emit(T_COMMENT)
			return aslInitial

		case eofRune:
			l.errorf("Unterminated comment")
			return nil
		}
	}
}

func aslSymbol(l *lexer) stateFn {
	for {
		r := l.accept(symbol, false)

		if r != utf8.RuneError && r != eofRune {
			s := l.current()

			if _, ok := symbolMap[s]; ok {
				continue
			}

			l.backup()
		}

		s := l.current()

		if t, ok := symbolMap[s]; ok {
			if t == T_COMMENT {
				return aslComment
			}

			if t == T_CLOSE_BRACE {
				l.emitGhost(T_TERMINATOR)
			}

			l.emit(t)
			return aslInitial
		}

		switch l.accept(symbol, false) {
		case eofRune, utf8.RuneError:
			l.errorf("Unknown symbol %q", s)
			return nil
		}
	}

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

func aslVariableName(l *lexer) stateFn {
	r := l.accept(identifier, true)

	switch r {
	case eofRune, utf8.RuneError:
		l.errorf("Invalid variable name `$`")
		return nil

	default:
		l.acceptRun(identifier, true)
		l.emit(T_VARIABLE)

		return aslInitial
	}
}
