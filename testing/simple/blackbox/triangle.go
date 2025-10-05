package simple

import (
	"fmt"
	"math"
)

func CalcArea(a, b, c uint) (float64, error) {
	if !isValidTriangle(a, b, c) {
		return 0, fmt.Errorf("the triagnle with sides %d, %d, %d is invalid", a, b, c)
	}
	sp := calcSemiperimeter(a, b, c)
	return math.Sqrt(sp * (sp - float64(a)) * (sp - float64(b)) * (sp - float64(c))), nil
}

func calcSemiperimeter(a, b, c uint) float64 {
	return float64((a + b + c)) / 2
}

func isValidTriangle(a, b, c uint) bool {
	return a+b > c && a+c > b && b+c > a
}
