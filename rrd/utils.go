package rrd

import (
	"unicode"
	"strings"
	"math"
	"strconv"
)

func rrdIsUnsignedInt(str string) bool {
	for _, ch := range str {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func rrdIsSignedInt(str string) bool {
	for i, ch := range str {
		if i == 0 && ch == '-' {
			continue
		}
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true

}

func rrdDiff(a, b string) float64 {
	a = strings.TrimFunc(a, func(ch rune) bool { return !unicode.IsDigit(ch)})
	aNeg := false
	if a[0] == '-' {
		aNeg = true
		a = a[1:]
	}

	b = strings.TrimFunc(b, func(ch rune) bool { return !unicode.IsDigit(ch)})
	bNeg := false
	if b[0] == '-' {
		bNeg = true
		b = b[1:]
	}
	result := math.NaN()
	m := len(a)
	if m < len(b) {
		m = len(b)
	}
	c := byte(0)
	r := make([]byte, m)
	aIdx := len(a) - 1
	bIdx := len(b) - 1
	for x := m - 1; x >= 0; x-- {
		if aIdx >= 0  && bIdx >= 0 {
			r[x] = a[aIdx] - c - b[bIdx] + '0'
		} else if aIdx >= 0 {
			r[x] = a[aIdx] - c
		} else {
			r[x] = ('0' - b[bIdx] - c) + '0'
		}
		if r[x] < '0' {
			r[x] += 10
			c = 1
		} else if r[x] > '9' {
			r[x] -= 10
			c = 1
		} else {
			c = 0
		}
		aIdx--
		bIdx--
	}

	var err error
	if c > 0 {
		for x:= m - 1; x >= 0; x-- {
			r[x] =('9' - r[x] + c) + '0'
			if r[x] > '9' {
				r[x] -= 10
				c = 1
			} else {
				c = 0
			}
		}

		result, err = strconv.ParseFloat(string(r), 64)
		if err != nil {
			result = math.NaN()
		} else {
			result = -result
		}
	} else {
		result, err = strconv.ParseFloat(string(r), 64)
		if err != nil {
			result = math.NaN()
		}
	}

	if aNeg && bNeg {
		result = -result
	}

	return result
}