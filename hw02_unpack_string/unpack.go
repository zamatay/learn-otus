package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type Data struct {
	prevRune      rune
	printString   string
	stringBuilder strings.Builder
}

func (r *Data) addItem(count int) {
	if r.printString == "" {
		return
	}
	t := strings.Repeat(r.printString, count)
	r.stringBuilder.WriteString(t)
	r.clearResult()
}

func (r *Data) clearResult() {
	r.printString, r.prevRune = "", 0
}

func (r *Data) Add(ru rune) {
	count, _ := strconv.Atoi(string(ru))
	r.addItem(count)
}

func (r *Data) setPrintRune(ru rune) {
	r.prevRune = 0
	r.printString = string(ru)
}

func isDigit(r rune) bool {
	return r >= 48 && r <= 57
}

func isSlash(r rune) bool {
	return r == '\\'
}

func isPrint(r rune) bool {
	return !isSlash(r) && !isDigit(r)
}

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	var result Data
	for _, currentRune := range s {
		if !result.getNextItem(currentRune) {
			return "", ErrInvalidString
		}
	}
	result.addItem(1)
	return result.stringBuilder.String(), nil
}

func (r *Data) getNextItem(currentRune int32) bool {
	switch {
	case isSlash(currentRune) && !isSlash(r.prevRune):
	case !isPrint(currentRune) && isSlash(r.prevRune):
		if isDigit(currentRune) && (!isSlash(r.prevRune)) {
			r.Add(currentRune)
		} else {
			r.addItem(1)
			r.setPrintRune(currentRune)
		}
		return true
	case isDigit(currentRune) && r.printString != "":
		r.Add(currentRune)
	case isPrint(currentRune):
		r.addItem(1)
		r.setPrintRune(currentRune)
	default:
		return false
	}
	r.prevRune = currentRune
	return true
}
