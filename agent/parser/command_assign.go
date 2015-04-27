package parser

type assignCommand struct {
	variable token
	expr     expression
}

func newAssignCommand(variable token, expr expression) command {
	return &assignCommand{variable, expr}
}

func (a *assignCommand) execute(c *executionContext) error {
	if vv, ok := a.expr.(resolvable); ok {
		vvv, err := vv.resolve(c)

		if err != nil {
			return err
		}

		res, err := expressionFromInterface(vvv, a.variable.line, a.variable.start)

		if err != nil {
			return err
		}

		c.variables[a.variable.source] = res
		return nil
	} else if val, err := a.expr.evaluate(c); err == nil {
		c.variables[a.variable.source] = val
		return nil
	} else {
		return err
	}
}
