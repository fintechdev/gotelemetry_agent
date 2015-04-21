//go:generate go tool yacc -o internal_parser.go -p parser internal_parser.y
package parser

func Parse(name, source string) ([]command, []error) {
	lexer := newASLLexer(lexASL(name, source))

	parserParse(lexer)

	return lexer.commands, lexer.errors
}
