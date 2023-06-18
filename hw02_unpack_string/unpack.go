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
		if !getNextItem(currentRune, &result) {
			return "", ErrInvalidString
		}
	}
	addItem(1, &result)
	return result.stringBuilder.String(), nil
}

func getNextItem(currentRune int32, result *Data) bool {
	switch {
	case isSlash(currentRune) && !isSlash(result.prevRune):
	case (isSlash(currentRune) || isDigit(currentRune)) && isSlash(result.prevRune):
		count := 1
		if isDigit(currentRune) && !isSlash(result.prevRune) {
			count, _ = strconv.Atoi(string(currentRune))
		}
		addItem(count, result)
		result.printString = string(currentRune)
		return true
	case isDigit(currentRune) && result.printString != "":
		count, _ := strconv.Atoi(string(currentRune))
		addItem(count, result)
	case isPrint(currentRune):
		addItem(1, result)
		result.printString = string(currentRune)
	default:
		return false
	}
	result.prevRune = currentRune
	return true
}

func addItem(count int, d *Data) {
	if d.printString == "" {
		return
	}
	t := strings.Repeat(d.printString, count)
	d.stringBuilder.WriteString(t)
	clearResult(d)
}

func clearResult(d *Data) {
	d.printString, d.prevRune = "", 0
}
