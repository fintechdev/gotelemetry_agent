package parser

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"regexp"
	"strings"
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

var excelExpressionRangeRegex = regexp.MustCompile(`^\s*([A-Za-z]+)(\d+)\s*[-:]\s*([A-Za-z]+)(\d+)\s*`)
var excelExpressionCellRegex = regexp.MustCompile(`^\s*([A-Za-z]+)(\d+)\s*$`)

func (x *excelExpression) parseRange(rangeSpec string) ([]excelExpressionCellReference, error) {
	result := []excelExpressionCellReference{}
	ranges := strings.Split(rangeSpec, ",")

	for _, r := range ranges {
		rangeMatches := excelExpressionRangeRegex.FindStringSubmatch(r)
		cellMatches := excelExpressionCellRegex.FindStringSubmatch(r)

		switch {
		case len(rangeMatches) > 0:
			startCell := newExcelExpressionCellReference(rangeMatches[1], rangeMatches[2])
			endCell := newExcelExpressionCellReference(rangeMatches[3], rangeMatches[4])

			if startCell.Column-endCell.Column != 0 && startCell.Row-endCell.Row != 0 {
				return result, errors.New("The range expression `" + r + "` does not represent a mono-dimensional block of cells.")
			}

			if startCell.Column-endCell.Column == 0 {
				var start, end int

				if startCell.Row < endCell.Row {
					start = startCell.Row - 1
					end = endCell.Row
				} else {
					start = endCell.Row - 1
					end = startCell.Row
				}

				for index := start; index < end; index++ {
					result = append(result, excelExpressionCellReference{Row: index, Column: startCell.Column})
				}
			} else {
				var start, end int

				if startCell.Column < endCell.Column {
					start = startCell.Column
					end = endCell.Column + 1
				} else {
					start = endCell.Column
					end = startCell.Column + 1
				}

				for index := start; index < end; index++ {
					result = append(result, excelExpressionCellReference{Row: startCell.Row, Column: index})
				}
			}

		case len(cellMatches) > 0:
			result = append(result, newExcelExpressionCellReference(cellMatches[1], cellMatches[2]))

		default:
			return result, errors.New("Unable to parse range expression `" + r + "`")
		}
	}

	return result, nil
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

			cells, err := x.parseRange(args["ranges"].(string))

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
					return nil, errors.New(fmt.Sprintf("Unable to handle value of type %s", cell.Type()))
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
