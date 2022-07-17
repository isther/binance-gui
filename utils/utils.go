package utils

import (
	"fmt"
	"math"
	"strings"
)

func Float64ToStringLen3(f float64) string {
	var s = fmt.Sprintf("%.8f", f)
	if f >= 1.0 {
		if f >= 100.0 {
			s = s[:3]
		} else {
			s = s[:4]
		}
	} else {
		for i := 0; i < len(s); i++ {
			if s[i] == '.' {
				continue
			}

			if s[i] != '0' {
				if i+3 <= len(s) {
					s = s[:i+3]
				}
				break
			}

		}

		var pos = len(s) - 1
		for ; pos > 0; pos-- {
			if s[pos] != '0' {
				break
			}
		}
		s = s[:pos+1]
	}
	if s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}

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
