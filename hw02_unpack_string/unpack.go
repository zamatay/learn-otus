package hw02unpackstring

import (
	"errors"
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

func (r *Data) checkItemIsNotValid(currentRune rune) bool {
	return (isDigit(r.prev2Rune) && isDigit(r.prevRune)) || (!isSlash(r.prev2Rune) && isDigit(currentRune) && isDigit(r.prevRune))
}
func (r *Data) getNextItem(currentRune int32) bool {
	if r.checkItemIsNotValid(currentRune) {
		return false
	}
	return true
}
