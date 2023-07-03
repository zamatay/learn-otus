package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

type Data struct {
	prevRune      rune
	prev2Rune     rune
	stringBuilder strings.Builder
}

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
	r.leftShift()
	r.printValue(0)
	r.printValue(0)
	return r.stringBuilder.String(), nil
}

func parseValue(countRune interface{}) int {
	switch v := countRune.(type) {
	case int:
		return v
	case rune:
		count, _ := strconv.Atoi(string(v))
		return count
	case string:
		count, _ := strconv.Atoi(v)
		return count
	default:
		return 0
	}
}

func isDigit(r rune) bool {
	return r >= 48 && r <= 57
}

func isSlash(r rune) bool {
	return r == '\\'
}

func isPrint(r rune) bool {
	return !isSlash(r) && !isDigit(r) && r != 0
}

func (r *Data) setPrevRune(currentRune rune) {
	r.prev2Rune = r.prevRune
	r.prevRune = currentRune
}

func (r *Data) checkItem(currentRune rune) bool {
	result := isDigit(r.prev2Rune) && isDigit(r.prevRune) || !isSlash(r.prev2Rune) && isDigit(currentRune) && isDigit(r.prevRune)
	return result
}

func (r *Data) addItem(item, currentRune rune, count interface{}, countShift int) {
	if item != 0 {
		cnt := parseValue(count)
		t := strings.Repeat(string(item), cnt)
		r.stringBuilder.WriteString(t)
	}

	switch countShift {
	case 1:
		r.prev2Rune, r.prevRune = r.prevRune, currentRune
	case 2:
		r.prev2Rune, r.prevRune = 0, currentRune
	case 3:
		r.prev2Rune, r.prevRune = 0, 0
	}
}

func (r *Data) getNextItem(currentRune int32) bool {
	if r.checkItem(currentRune) {
		return false
	}
	if r.prevRune > 0 && r.prev2Rune > 0 {
		r.printValue(currentRune)
	} else {
		r.setPrevRune(currentRune)
	}
	return true
}

func (r *Data) getItems(currentRune rune) (rune, rune, rune) {
	switch {
	case r.prev2Rune != 0:
		return r.prev2Rune, r.prevRune, currentRune
	case r.prevRune != 0:
		return r.prevRune, currentRune, 0
	default:
		return currentRune, 0, 0
	}
}

func (r *Data) printValue(currentRune rune) {
	item1, item2, item3 := r.getItems(currentRune)
	switch {
	case isPrint(item1) && !isDigit(item2):
		r.addItem(item1, currentRune, 1, 1)
	case isPrint(item1) && isDigit(item2):
		r.addItem(item1, currentRune, item2, 2)
	case isSlash(item1) && !isPrint(item2) && isDigit(item3):
		r.addItem(item2, currentRune, item3, 3)
	case isSlash(item1) && !isPrint(item2):
		r.addItem(item2, currentRune, 1, 2)
	}
}

func (r *Data) leftShift() {
	r.prev2Rune, r.prevRune, _ = r.getItems(0)
}
