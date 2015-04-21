package parser

type assignCommand struct {
	variable token
	expr     expression
}

func newAssignCommand(variable token, expr expression) command {
	return &assignCommand{variable, expr}
}

func (a *assignCommand) execute(c *executionContext) error {
	if val, err := a.expr.evaluate(c); err == nil {
		c.variables[a.variable.source] = val
		return nil
	} else {
		return err
	}
}
