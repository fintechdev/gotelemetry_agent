package parser

type expression interface {
	evaluate(c *executionContext) (interface{}, error)
	line() int
	position() int
}
