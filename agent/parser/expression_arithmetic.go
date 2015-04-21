package parser

import (
	"errors"
)

func forceNumeric(c *executionContext, left, right expression, line, position int) (float64, float64, error) {
	l, err := left.evaluate(c)

	if err != nil {
		return 0, 0, err
	}

	r, err := right.evaluate(c)

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

type artihmeticExpression struct {
	left     expression
	right    expression
	operator terminal
	l        int
	p        int
}

func newArithmeticExpression(left, right expression, operator terminal, line, position int) expression {
	return &artihmeticExpression{left, right, operator, line, position}
}

func (a *artihmeticExpression) evaluate(c *executionContext) (interface{}, error) {
	l, r, err := forceNumeric(c, a.left, a.right, a.l, a.p)

	switch a.operator {
	case T_PLUS:
		return l + r, err

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

func (a *artihmeticExpression) line() int {
	return a.l
}

func (a *artihmeticExpression) position() int {
	return a.p
}
