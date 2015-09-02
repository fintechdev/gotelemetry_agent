package parser

type Command interface {
	execute(c *executionContext) error
}
