package parser

type evaluateCommand struct {
	expr expression
}

func newEvaluateCommand(expr expression) Command {
	return &evaluateCommand{expr}
}

func (e *evaluateCommand) execute(c *executionContext) error {
	_, err := resolveExpression(c, e.expr)
	return err
}
