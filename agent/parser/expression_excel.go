package parser

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
)

// Number

type excelExpression struct {
	path string
	xl   *xlsx.File
	err  error
	l    int
	p    int
}

func newExcelExpression(path string, line, position int) expression {
	xl, err := xlsx.OpenFile(path)

	result := &excelExpression{
		xl:   xl,
		err:  err,
		path: path,
		l:    line,
		p:    position,
	}

	return result
}

func (x *excelExpression) evaluate(c *executionContext) (interface{}, error) {
	return nil, errors.New(fmt.Sprintf("Cannot evaluate an expression of type %s", x))
}

// Properties

type excelProperty func(x *excelExpression) expression

var excelProperties = map[string]excelProperty{
	"cells": func(x *excelExpression) expression {
		return x.cells()
	},
}

func (x *excelExpression) cells() expression {
	return newCallableExpression(
		"cells",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			if x.err != nil {
				return nil, x.err
			}

			var sheet *xlsx.Sheet
			var sheets = x.xl.Sheets

			if len(sheets) == 0 {
				return nil, errors.New(fmt.Sprintf("The spreadsheet at %s is empty.", x.path))
			}

			if sheetIndex, ok := args["sheet"].(int); ok {
				if sheetIndex < 0 || sheetIndex >= len(sheets) {
					return nil, errors.New(fmt.Sprintf("Invalid sheet index %d (valid values are 0–%d)", sheetIndex, len(sheets)))
				}

				sheet = sheets[sheetIndex]
			} else {
				sheet = sheets[sheetIndex]
			}

			cells, err := parseRange(args["ranges"].(string))

			if err != nil {
				return nil, err
			}

			result := make([]interface{}, len(cells))

			for index, position := range cells {
				cell := sheet.Cell(position.Row, position.Column)

				switch cell.Type() {
				case xlsx.CellTypeBool:
					result[index] = newBooleanExpression(cell.Bool(), x.l, x.p)

				case xlsx.CellTypeNumeric:
					if v, err := cell.Float(); err == nil {
						result[index] = newNumericExpression(v, x.l, x.p)
					} else {
						return nil, err
					}

				case xlsx.CellTypeString:
					result[index] = newStringExpression(cell.String(), x.l, x.p)

				default:
					return nil, errors.New(fmt.Sprintf("Unable to handle value of type %d", cell.Type()))
				}
			}

			return newArrayExpression(result, x.l, x.p), nil
		},
		map[string]callableArgument{
			"sheet":  callableArgumentOptionalNumeric,
			"ranges": callableArgumentString,
		},
		x.l,
		x.p,
	)
}

func (x *excelExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := excelProperties[property]; ok {
		return f(x), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", x, property))
}

func (x *excelExpression) line() int {
	return x.l
}

func (x *excelExpression) position() int {
	return x.p
}

func (x *excelExpression) String() string {
	return fmt.Sprintf("Excel(%v)", x.path)
}
