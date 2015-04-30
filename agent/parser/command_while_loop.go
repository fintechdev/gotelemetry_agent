package parser

type whileLoopCommand struct {
	condition expression
	commands  []command
}

func newWhileLoopCommand(condition expression, commands []command) command {
	return &whileLoopCommand{condition, commands}
}

func (w *whileLoopCommand) execute(c *executionContext) error {
	for {
		expr, err := resolveExpression(c, w.condition)

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

		if !val {
			return nil
		}

		for _, cmd := range w.commands {
			if err := cmd.execute(c); err != nil {
				return err
			}
		}
	}

	return nil
}
