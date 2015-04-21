package parser

import (
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
	return p.target.extract(c, p.name)
}

func (p *propertyExpression) extract(c *executionContext, property string) (expression, error) {
	return p.target.extract(c, p.name)
}

func (p *propertyExpression) call(c *executionContext, arguments map[string]interface{}) (expression, error) {
	target, err := p.target.extract(c, p.name)

	if err != nil {
		return nil, err
	}

	return target.call(c, arguments)
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
