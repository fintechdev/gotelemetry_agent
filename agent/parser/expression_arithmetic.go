package parser

import (
	"errors"
	"fmt"
)

func forceNumeric(c *executionContext, left, right expression, line, position int) (float64, float64, error) {
	l, err := resolveExpression(c, left)

	if err != nil {
		return 0, 0, err
	}

	r, err := resolveExpression(c, right)

	if err != nil {
		return 0, 0, err
	}

	ll := newNumericExpression(l, line, position)
	rr := newNumericExpression(r, line, position)

	lll, err := ll.evaluate(c)

	if err != nil {
		return 0, 0, err
	}

	rrr, err := rr.evaluate(c)

	if err != nil {
		return 0, 0, err
	}

	return lll.(float64), rrr.(float64), nil
}

type arithmeticExpression struct {
	left     expression
	right    expression
	operator token
	l        int
	p        int
}

func newArithmeticExpression(left, right expression, operator token, line, position int) expression {
	return &arithmeticExpression{left, right, operator, line, position}
}

func (a *arithmeticExpression) evaluateAddition(c *executionContext) (interface{}, error) {
	l, r, err := forceNumeric(c, a.left, a.right, a.l, a.p)

	if err == nil {
		return l + r, nil
	}

	ll, err := resolveExpression(c, a.left)

	if err != nil {
		return nil, err
	}

	if lll, ok := ll.(string); ok {
		rrr, err := resolveExpression(c, a.right)

		if err != nil {
			return nil, err
		}

		return fmt.Sprintf("%s%v", lll, rrr), nil
	}

	return nil, errors.New(fmt.Sprintf("Unable to add expressions %s and %s", a.left, a.right))
}

func (a *arithmeticExpression) evaluate(c *executionContext) (interface{}, error) {
	if a.operator.terminal == T_PLUS {
		return a.evaluateAddition(c)
	}

	l, r, err := forceNumeric(c, a.left, a.right, a.l, a.p)

	switch a.operator.terminal {
	case T_MULTIPLY:
		return l * r, err

	case T_MINUS:
		return l - r, err

	case T_DIVIDE:
		if err != nil {
			return 0, err
		}

		if r == 0.0 {
			return 0, errors.New("Divide by zero")
		}

		return l / r, err

	default:
		panic("Unknown operator " + a.operator.String())
		return 0, nil
	}
}

func (a *arithmeticExpression) line() int {
	return a.l
}

func (a *arithmeticExpression) position() int {
	return a.p
}

func (a *arithmeticExpression) String() string {
	return fmt.Sprintf("%s %s %s", a.left, a.operator.source, a.right)
}
