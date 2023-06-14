package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type ResultData struct {
	s               string
	countRepeat     int
	currentPosition int
}

func isDigit(r uint8) bool {
	return r >= 48 && r <= 57
}

func isSlash(r uint8) bool {
	return r == 92
}

func NewResultData(item, cnt uint8, position int) *ResultData {
	countRepeat, _ := strconv.Atoi(string(cnt))

	return &ResultData{
		s:               string(item),
		countRepeat:     countRepeat,
		currentPosition: position,
	}
}

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if isDigit(s[0]) {
		return "", ErrInvalidString
	}
	var b strings.Builder
	currentIndex := 0
	for {
		resultData, err := getNextItem(s, currentIndex)
		if err != nil {
			return "", err
		}
		currentIndex = resultData.currentPosition
		b.WriteString(strings.Repeat(resultData.s, resultData.countRepeat))
		if len(s)-1 < currentIndex {
			break
		}
	}
	println(b.String())

	return b.String(), nil
}

func getNextItem(s string, index int) (*ResultData, error) {
	var item1, item2, item3 uint8
	item1 = s[index]
	if len(s)-1 >= index+1 {
		item2 = s[index+1]
	}
	if len(s)-1 >= index+2 {
		item3 = s[index+2]
	}
	if isDigit(item1) && isDigit(item2) || !isSlash(item1) && isDigit(item2) && isDigit(item3) {
		return nil, ErrInvalidString
	}
	switch {
	case isSlash(item1) && isDigit(item3):
		return NewResultData(item2, item3, index+3), nil
	case isSlash(item1) && isSlash(item2):
		return NewResultData(item1, '1', index+2), nil
	case isSlash(item1) && isDigit(item2):
		if item3 != 0 && isDigit(item3) {
			return NewResultData(item2, item3, index+3), nil
		}
		return NewResultData(item2, '1', index+2), nil
	case !isDigit(item1) && isDigit(item2):
		return NewResultData(item1, item2, index+2), nil
	default:
		return NewResultData(item1, '1', index+1), nil
	}
}
