package parser

import (
	"fmt"
)

func forceBoolean(c *executionContext, left, right expression, line, position int) (bool, bool, error) {
	l, err := resolveExpression(c, left)

	if err != nil {
		return false, true, err
	}

	r, err := resolveExpression(c, right)

	if err != nil {
		return false, true, err
	}

	ll := newBooleanExpression(l, line, position)
	rr := newBooleanExpression(r, line, position)

	lll, err := ll.evaluate(c)

	if err != nil {
		return false, true, err
	}

	rrr, err := rr.evaluate(c)

	if err != nil {
		return false, true, err
	}

	return lll.(bool), rrr.(bool), nil
}

func forceComparable(c *executionContext, left, right expression, line, position int) (interface{}, interface{}, error) {
	l, err := resolveExpression(c, left)

	if err != nil {
		return 0, 0, err
	}

	r, err := resolveExpression(c, right)

	if err != nil {
		return 0, 0, err
	}

	if _, ok := l.(bool); ok {
		if _, ok := r.(bool); ok {
			return l, r, nil
		}
	}

	ll := newNumericExpression(l, line, position)
	rr := newNumericExpression(r, line, position)

	lll, err1 := ll.evaluate(c)
	rrr, err2 := rr.evaluate(c)

	if err1 == nil && err2 == nil {
		return lll, rrr, nil
	}

	return l, r, nil
}

type logicalExpression struct {
	left     expression
	right    expression
	operator token
	l        int
	p        int
}

func newLogicalExpression(left, right expression, operator token, line, position int) expression {
	return &logicalExpression{left, right, operator, line, position}
}

func (x *logicalExpression) evaluate(c *executionContext) (interface{}, error) {
	switch x.operator.terminal {
	case T_EQUAL:
		l, r, err := forceComparable(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l == r, x.l, x.p), nil

	case T_NOT_EQUAL:
		l, r, err := forceComparable(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l != r, x.l, x.p), nil

	case T_GREATER_THAN:
		l, r, err := forceNumeric(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l > r, x.l, x.p), nil

	case T_GREATER_THAN_OR_EQUAL:
		l, r, err := forceNumeric(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l >= r, x.l, x.p), nil

	case T_LESS_THAN:
		l, r, err := forceNumeric(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l < r, x.l, x.p), nil

	case T_LESS_THAN_OR_EQUAL:
		l, r, err := forceNumeric(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l <= r, x.l, x.p), nil

	case T_AND:
		l, r, err := forceBoolean(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l && r, x.l, x.p), nil

	case T_OR:
		l, r, err := forceBoolean(c, x.left, x.right, x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(l || r, x.l, x.p), nil

	case T_NEGATE:
		l, _, err := forceBoolean(c, x.left, newBooleanExpression(false, 0, 0), x.l, x.p)

		if err != nil {
			return nil, err
		}

		return newBooleanExpression(!l, x.l, x.p), nil

	default:
		panic("Unknown operator " + x.operator.String())
		return 0, nil
	}
}

func (x *logicalExpression) line() int {
	return x.l
}

func (x *logicalExpression) position() int {
	return x.p
}

func (x *logicalExpression) String() string {
	return fmt.Sprintf("%s %s %s", x.left, x.operator.source, x.right)
}
