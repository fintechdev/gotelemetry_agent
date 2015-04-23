package parser

type ifThenElseCommand struct {
	condition    expression
	thenCommands []command
	elseCommands []command
}

func newIfThenElseCommand(condition expression, thenCommands, elseCommands []command) command {
	return &ifThenElseCommand{condition, thenCommands, elseCommands}
}

func (i *ifThenElseCommand) execute(c *executionContext) error {
	expr, err := resolveExpression(c, i.condition)

	if err != nil {
		return err
	}

	val, ok := expr.(bool)

	if !ok {
		cond := newBooleanExpression(expr, 0, 0)

		v, err := cond.evaluate(c)

		if err != nil {
			return err
		}

		val = v.(bool)
	}

	commands := i.thenCommands

	if !val {
		commands = i.elseCommands
	}

	for _, cmd := range commands {
		if err := cmd.execute(c); err != nil {
			return err
		}
	}

	return nil
}
