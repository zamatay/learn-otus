package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type Data struct {
	prevRune      rune
	prev2Rune     rune
	stringBuilder strings.Builder
}

func Unpack(_ string) (string, error) {
	return "", nil
}
