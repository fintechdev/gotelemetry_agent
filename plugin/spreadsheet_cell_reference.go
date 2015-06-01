package plugin

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type spreadsheetCellReference struct {
	Row    int
	Column int
}

func newExcelCellReference(column, row string) spreadsheetCellReference {
	result := spreadsheetCellReference{}

	column = strings.ToUpper(column)

	l := len(column) - 1

	for index, char := range column {
		result.Column += int(float64((char-'A')+1)*(math.Pow(26, float64(l-index)))) - 1
	}

	var err error

	result.Row, err = strconv.Atoi(row)

	if err != nil {
		panic(err) // This should never happen, since we always come here from a regex that only matches digits
	}

	return result
}

var spreadsheetRangeRegex = regexp.MustCompile(`^\s*([A-Za-z]+)(\d+)\s*[-:]\s*([A-Za-z]+)(\d+)\s*`)
var spreadsheetCellRegex = regexp.MustCompile(`^\s*([A-Za-z]+)(\d+)\s*$`)

func parseRange(rangeSpec string) ([]spreadsheetCellReference, error) {
	result := []spreadsheetCellReference{}
	ranges := strings.Split(rangeSpec, ",")

	for _, r := range ranges {
		rangeMatches := spreadsheetRangeRegex.FindStringSubmatch(r)
		cellMatches := spreadsheetCellRegex.FindStringSubmatch(r)

		switch {
		case len(rangeMatches) > 0:
			startCell := newExcelCellReference(rangeMatches[1], rangeMatches[2])
			endCell := newExcelCellReference(rangeMatches[3], rangeMatches[4])

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
					result = append(result, spreadsheetCellReference{Row: index, Column: startCell.Column})
				}
			} else {
				var start, end int

				if startCell.Column < endCell.Column {
					start = startCell.Column - 1
					end = endCell.Column
				} else {
					start = endCell.Column - 1
					end = startCell.Column
				}

				for index := start; index < end; index++ {
					result = append(result, spreadsheetCellReference{Row: startCell.Row, Column: index})
				}
			}

		case len(cellMatches) > 0:
			result = append(result, newExcelCellReference(cellMatches[1], cellMatches[2]))

		default:
			return result, errors.New("Unable to parse range expression `" + r + "`")
		}
	}

	return result, nil
}
