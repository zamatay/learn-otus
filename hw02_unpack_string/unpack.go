package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type ResultData struct {
	printRune       rune
	prevRune        rune
	countRepeat     int
	currentPosition int
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
	var result strings.Builder
	item := ResultData{}
	item.currentPosition = -1
	isPrintDigit := false
	for i, currentRune := range s {
		item.currentPosition++
		if (i == 0 || isPrintDigit) && isDigit(currentRune) {
			return "", ErrInvalidString
		}

		if item.currentPosition == 0 && isSlash(currentRune) {
			item.prevRune = currentRune
			continue
		}
		isPrintDigit = false
		switch {
		case item.printRune != 0 && isDigit(currentRune) && !isSlash(item.prevRune):
			item.countRepeat, _ = strconv.Atoi(string(currentRune))
			isPrintDigit = true
			addItem(&result, &item)
		case isPrint(currentRune) || isSlash(item.prevRune):
			addItem(&result, &item)
			item.printRune = currentRune
			item.countRepeat = 1
		}
	}
	addItem(&result, &item)
	return result.String(), nil
}

func addItem(s *strings.Builder, r *ResultData) {
	if r.printRune == 0 {
		return
	}
	t := strings.Repeat(string(r.printRune), r.countRepeat)
	s.WriteString(t)
	clearResult(r)
}

func clearResult(r *ResultData) {
	r.currentPosition = -1
	r.countRepeat = 1
	r.printRune = 0
	r.prevRune = 0
}
