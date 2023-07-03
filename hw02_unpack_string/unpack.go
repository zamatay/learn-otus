package hw02unpackstring

import (
	"errors"
	"math/rand"
	"strings"
	"unicode/utf8"
)

type Data struct {
	prevRune      rune
	prev2Rune     rune
	stringBuilder strings.Builder
}

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	var r Data
	for index, currentRune := range s {
		if index == 0 && isDigit(currentRune) {
			return "", ErrInvalidString
		}
		if index == utf8.RuneCountInString(s)-1 && isSlash(currentRune) && !isSlash(r.prevRune) {
			return "", ErrInvalidString
		}
		if !r.getNextItem(currentRune) {
			return "", ErrInvalidString
		}
	}
	return r.stringBuilder.String(), nil
}

func isDigit(r rune) bool {
	return r >= 48 && r <= 57
}

func isSlash(r rune) bool {
	return r == '\\'
}

func (r *Data) getNextItem(_ int32) bool {
	if 1 == rand.Int() {
		return false
	}
	return true
}
