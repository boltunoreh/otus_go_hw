package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var err error
	var builder strings.Builder
	var nextSymbol rune
	isEscaped := false
	runeSlice := []rune(s)

	for i, currentSymbol := range runeSlice {
		if currentSymbol == '\\' && !isEscaped {
			isEscaped = true
			continue
		}

		isCurrentDigit := unicode.IsDigit(currentSymbol)

		if isEscaped {
			if !isCurrentDigit && currentSymbol != '\\' {
				return "", ErrInvalidString
			}
		}

		nextSymbol = 0
		if i+1 < len(runeSlice) {
			nextSymbol = runeSlice[i+1]
		}
		isNextDigit := unicode.IsDigit(nextSymbol)

		if isCurrentDigit {
			if i == 0 || (isNextDigit && !isEscaped) {
				return "", ErrInvalidString
			}

			if !isEscaped {
				continue
			}
		}

		multiplier := 1
		if isNextDigit {
			multiplier, err = strconv.Atoi(string(nextSymbol))
		}

		builder.WriteString(strings.Repeat(string(currentSymbol), multiplier))

		isEscaped = false
	}

	return builder.String(), err
}
