package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type googleSpreadsheetExpression struct {
	spreadsheetId string
	worksheetIds  []string
	err           error
	l             int
	p             int
}

type cellField struct {
	XMLName      xml.Name `xml:"cell"`
	NumericValue string   `xml:"numericValue,attr"`
	Value        string   `xml:",innerxml"`
}

type entryField struct {
	XMLName xml.Name  `xml:"entry"`
	Id      string    `xml:"id"`
	Cell    cellField `xml:"cell"`
}

type atomFeed struct {
	XMLName xml.Name     `xml:"feed"`
	Entries []entryField `xml:"entry"`
}

func newGoogleSpreadsheetExpression(spreadsheetId string, line, position int) expression {
	res, err := http.Get(fmt.Sprintf("https://spreadsheets.google.com/feeds/worksheets/%s/public/full", spreadsheetId))
	if err != nil {
		return &googleSpreadsheetExpression{err: err}
	}

	rawXML, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &googleSpreadsheetExpression{err: err}
	}

	spreadsheetFeed := atomFeed{}
	err = xml.Unmarshal(rawXML, &spreadsheetFeed)
	if err != nil {
		return &googleSpreadsheetExpression{err: errors.New(fmt.Sprintf("The Google spreadsheet with the ID `%s` is not published. Spreadsheets must be published publicly in order to be accessible", spreadsheetId))}
	}

	worksheetIds := []string{}

	for _, entry := range spreadsheetFeed.Entries {
		worksheetURL := entry.Id
		worksheetId := strings.Replace(worksheetURL, fmt.Sprintf("https://spreadsheets.google.com/feeds/worksheets/%s/public/full/", spreadsheetId), "", 1)
		worksheetIds = append(worksheetIds, worksheetId)
	}

	result := &googleSpreadsheetExpression{
		spreadsheetId: spreadsheetId,
		worksheetIds:  worksheetIds,
		err:           nil,
		l:             line,
		p:             position,
	}

	return result
}

func (x *googleSpreadsheetExpression) evaluate(c *executionContext) (interface{}, error) {
	return nil, errors.New(fmt.Sprintf("Cannot evaluate an expression of type %s", x))
}

// Properties

type googleSpreadsheetProperty func(x *googleSpreadsheetExpression) expression

var googleSpreadsheetProperties = map[string]googleSpreadsheetProperty{
	"cells": func(x *googleSpreadsheetExpression) expression {
		return x.cells()
	},
}

func (x *googleSpreadsheetExpression) cells() expression {
	return newCallableExpression(
		"cells",
		func(c *executionContext, args map[string]interface{}) (expression, error) {
			if x.err != nil {
				return nil, x.err
			}

			var worksheetId string

			if len(x.worksheetIds) == 0 {
				return nil, errors.New(fmt.Sprintf("The spreadsheet %s is empty.", x.spreadsheetId))
			}

			if sheetIndex, ok := args["sheet"].(int); ok {
				if sheetIndex < 0 || sheetIndex >= len(x.worksheetIds) {
					return nil, errors.New(fmt.Sprintf("Invalid sheet index %d (valid values are 0–%d)", sheetIndex, len(x.worksheetIds)))
				}
				worksheetId = x.worksheetIds[sheetIndex]
			} else {
				worksheetId = x.worksheetIds[0]
			}

			cells, err := parseRange(args["ranges"].(string))

			if err != nil {
				return nil, err
			}

			result := make([]interface{}, len(cells))

			for index, position := range cells {
				cellURL := fmt.Sprintf("https://spreadsheets.google.com/feeds/cells/%s/%s/public/full?min-row=%d&max-row=%d&min-col=%d&max-col=%d ", x.spreadsheetId, worksheetId, position.Row-1, position.Row-1, position.Column+1, position.Column+1)

				res, err := http.Get(cellURL)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Unable to get cell %s: %v", cellURL, err))
				}

				rawXML, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Unable to read cell data %s: %v", cellURL, err))

				}

				cellFeed := atomFeed{}
				err = xml.Unmarshal(rawXML, &cellFeed)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Unable to unmarshal cell data %s: %v", cellURL, err))
				}

				c := cellFeed.Entries[0].Cell

				if c.Value == "FALSE" || c.Value == "TRUE" {
					result[index] = newBooleanExpression((c.Value == "TRUE"), x.l, x.p)
				} else if len(c.NumericValue) > 0 {
					numericValue, err := strconv.ParseFloat(c.NumericValue, 64)
					if err != nil {
						return nil, errors.New(fmt.Sprintf("Unable to parse numeric value of cell data %s: %v", cellURL, err))
					}
					result[index] = newNumericExpression(numericValue, x.l, x.p)
				} else {
					result[index] = newStringExpression(c.Value, x.l, x.p)
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

func (x *googleSpreadsheetExpression) extract(c *executionContext, property string) (expression, error) {
	if f, ok := googleSpreadsheetProperties[property]; ok {
		return f(x), nil
	}

	return nil, errors.New(fmt.Sprintf("%s does not contain a property with the key `%s`", x, property))
}

func (x *googleSpreadsheetExpression) line() int {
	return x.l
}

func (x *googleSpreadsheetExpression) position() int {
	return x.p
}

func (x *googleSpreadsheetExpression) String() string {
	return fmt.Sprintf("GoogleSpreadsheet(%v)", x.spreadsheetId)
}
