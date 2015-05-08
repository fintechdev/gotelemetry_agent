package parser

import (
	"math"
	"strconv"
	"strings"
)

type excelExpressionCellReference struct {
	Row    int
	Column int
}

func newExcelExpressionCellReference(column, row string) excelExpressionCellReference {
	result := excelExpressionCellReference{}

	column = strings.ToUpper(column)

	l := len(column) - 1

	for index, char := range column {
		result.Column += int(float64((char - 'A')) * (math.Pow(26, float64(l-index))))
	}

	var err error

	result.Row, err = strconv.Atoi(row)
	result.Row += 1

	if err != nil {
		panic(err) // This should never happen, since we always come here from a regex that only matches digits
	}

	return result
}
