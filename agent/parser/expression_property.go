package parser

import (
	"errors"
	"fmt"
)

// Number

type propertyExpression struct {
	target expression
	name   string
	l      int
	p      int
}

func newPropertyExpression(target expression, name string, line, position int) expression {
	result := &propertyExpression{
		target: target,
		name:   name,
		l:      line,
		p:      position,
	}

	return result
}

func (p *propertyExpression) evaluate(c *executionContext) (interface{}, error) {
	ex := p.target

	if exx, err := p.target.evaluate(c); err == nil {
		if _, ok := exx.(extractable); ok {
			ex = exx.(expression)
		}
	}

	if ex, ok := ex.(extractable); ok {
		return ex.extract(c, p.name)
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", p.target, p.name))
}

func (p *propertyExpression) extract(c *executionContext, property string) (expression, error) {
	ex := p.target

	if exx, err := p.target.evaluate(c); err == nil {
		if _, ok := exx.(extractable); ok {
			ex = exx.(expression)
		}
	}

	if ex, ok := ex.(extractable); ok {
		return ex.extract(c, p.name)
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", p.target, property))
}

func (p *propertyExpression) call(c *executionContext, arguments map[string]interface{}) (expression, error) {
	ex := p.target

	if exx, err := p.target.evaluate(c); err == nil {
		if _, ok := exx.(extractable); ok {
			ex = exx.(expression)
		}
	}

	if ex, ok := ex.(extractable); ok {
		target, err := ex.extract(c, p.name)

		if err != nil {
			return nil, err
		}

		if cl, ok := target.(callable); ok {
			return cl.call(c, arguments)
		}

		return nil, errors.New(fmt.Sprintf("%s is not a function", target))
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", p.target, p.name))
}

func (p *propertyExpression) line() int {
	return p.l
}

func (p *propertyExpression) position() int {
	return p.p
}

func (p *propertyExpression) String() string {
	return fmt.Sprintf("Property(%s.%s)", p.target, p.name)
}
