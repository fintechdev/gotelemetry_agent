package parser

type command interface {
	execute(c *executionContext) error
}
