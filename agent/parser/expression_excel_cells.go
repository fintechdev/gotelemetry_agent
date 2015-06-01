package parser

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type spreadsheetExpressionCellReference struct {
	Row    int
	Column int
}

func newSpreadsheetExpressionCellReference(column, row string) spreadsheetExpressionCellReference {
	result := spreadsheetExpressionCellReference{}

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

var spreadsheetExpressionRangeRegex = regexp.MustCompile(`^\s*([A-Za-z]+)(\d+)\s*[-:]\s*([A-Za-z]+)(\d+)\s*`)
var spreadsheetExpressionCellRegex = regexp.MustCompile(`^\s*([A-Za-z]+)(\d+)\s*$`)

func parseRange(rangeSpec string) ([]spreadsheetExpressionCellReference, error) {
	result := []spreadsheetExpressionCellReference{}
	ranges := strings.Split(rangeSpec, ",")

	for _, r := range ranges {
		rangeMatches := spreadsheetExpressionRangeRegex.FindStringSubmatch(r)
		cellMatches := spreadsheetExpressionCellRegex.FindStringSubmatch(r)

		switch {
		case len(rangeMatches) > 0:
			startCell := newSpreadsheetExpressionCellReference(rangeMatches[1], rangeMatches[2])
			endCell := newSpreadsheetExpressionCellReference(rangeMatches[3], rangeMatches[4])

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
					result = append(result, spreadsheetExpressionCellReference{Row: index, Column: startCell.Column})
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
					result = append(result, spreadsheetExpressionCellReference{Row: startCell.Row, Column: index})
				}
			}

		case len(cellMatches) > 0:
			result = append(result, newSpreadsheetExpressionCellReference(cellMatches[1], cellMatches[2]))

		default:
			return result, errors.New("Unable to parse range expression `" + r + "`")
		}
	}

	return result, nil
}
