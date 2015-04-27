package parser

import (
	"errors"
	"fmt"
	"math"
)

type numericArray []float64

func (a *arrayExpression) numericArray(c *executionContext) (numericArray, error) {
	result := make(numericArray, len(a.values))

	for index, value := range a.values {
		v, err := resolveExpression(c, newNumericExpression(value, a.l, a.p))

		if err != nil {
			return numericArray{}, errors.New(fmt.Sprintf("Cannot derive an array of numeric values, because converting the value at index %d caused this error: %s", index, err))
		}

		result[index] = v.(float64)
	}

	return result, nil
}

func (n numericArray) sum() float64 {
	result := 0.0

	for _, v := range n {
		result += v
	}

	return result
}

func (n numericArray) min() float64 {
	if len(n) == 0 {
		return 0.0
	}

	result := math.MaxFloat64

	for _, v := range n {
		result = math.Min(result, v)
	}

	return result
}

func (n numericArray) max() float64 {
	if len(n) == 0 {
		return 0.0
	}

	result := -math.MaxFloat64

	for _, v := range n {
		result = math.Max(result, v)
	}

	return result
}

func (n numericArray) avg() float64 {
	if len(n) == 0 {
		return 0.0
	}

	return n.sum() / float64(len(n))
}

func (n numericArray) stddev() float64 {
	avg := n.avg()
	squares := make(numericArray, len(n))

	for i, v := range n {
		squares[i] = (v - avg) * (v - avg)
	}

	return math.Sqrt(squares.avg())
}
