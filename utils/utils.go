package utils

import (
	"fmt"
	"math"
	"strings"
)

func correction(val float64, size string) string {
	var (
		oneIdx    = strings.Index(size, "1")
		pointIdx  = strings.Index(size, ".")
		precision int
		resStr    string
	)
	if oneIdx < pointIdx {
		precision = oneIdx - pointIdx + 1
		resStr = fmt.Sprintf("%.8f", RoundLower(val, precision))
		pointIdx = strings.Index(resStr, ".")
		resStr = resStr[:pointIdx]
		return resStr
	}

	precision = oneIdx - pointIdx
	resStr = fmt.Sprintf("%.8f", RoundLower(val, precision))

	return resStr
}

func RoundLower(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p) * math.Pow10(-precision)
	}

	return math.Floor(val*p) / p
}

func RoundUpper(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}
