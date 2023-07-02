package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type Data struct {
	prevRune      rune
	prev2Rune     rune
	stringBuilder strings.Builder
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

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	fmt.Println(s)
	var r Data
	for _, currentRune := range s {
		if !r.getNextItem(currentRune) {
			return "", ErrInvalidString
		}
	}
	r.printValue(0)
	r.printValue(0)
	return r.stringBuilder.String(), nil
}

func (r *Data) setPrevRune(currentRune rune) {
	r.prev2Rune = r.prevRune
	r.prevRune = currentRune
}

func (r *Data) checkLastItemIsNotValid() bool {
	return (isDigit(r.prevRune) && r.prev2Rune == 0) || (isDigit(r.prevRune) && isDigit(r.prev2Rune)) || (isSlash(r.prevRune) && isSlash(r.prev2Rune))
}

func (r *Data) addItem(item rune, count_rune interface{}) {
	if item == 0 {
		return
	}
	count := parseValue(count_rune)

	t := strings.Repeat(string(item), count)
	// fmt.Println(t)
	r.stringBuilder.WriteString(t)
}

func (r *Data) shift(item rune, isTwo bool) {
	if r.prev2Rune == 0 && item == 0 && r.prevRune != 0 {
		r.prev2Rune = r.prevRune
	}
	if isTwo {
		r.prev2Rune = 0
	} else {
		r.prev2Rune = r.prevRune
	}
	r.prevRune = item
}

func (r *Data) getNextItem(currentRune int32) bool {
	if r.prevRune > 0 && r.prev2Rune > 0 {
		r.printValue(currentRune)
	} else {
		r.setPrevRune(currentRune)
	}
	return true //!r.checkLastItemIsNotValid()
}

func (r *Data) getItems(currentRune rune) (rune, rune, rune) {
	if r.prev2Rune != 0 {
		return r.prev2Rune, r.prevRune, currentRune
	} else if r.prevRune != 0 {
		return r.prevRune, currentRune, 0
	} else {
		return currentRune, 0, 0
	}
}

func (r *Data) printValue(currentRune rune) {
	// если
	item1, item2, item3 := r.getItems(currentRune)
	if isPrint(item1) && !isDigit(item2) {
		r.addItem(item1, 1)
		r.shift(item3, false)
	} else if isPrint(item1) && isDigit(item2) {
		r.addItem(item1, item2)
		r.shift(item3, true)
	} else if isSlash(item1) && !isPrint(item2) && isDigit(item3) {
		r.addItem(item2, item3)
		r.shift(item3, true)
	} else if isSlash(item1) && isDigit(item2) {
		r.addItem(item2, 1)
		r.shift(item3, true)
	}

}
