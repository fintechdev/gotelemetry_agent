package parser

type outputCommand struct {
	property token
	expr     expression
}

func newOutputCommand(property token, expr expression) command {
	return &outputCommand{property, expr}
}

func (o *outputCommand) execute(c *executionContext) error {
	if val, err := resolveExpression(c, o.expr); err == nil {
		c.output[o.property.source] = val
		return nil
	} else {
		return err
	}
}
