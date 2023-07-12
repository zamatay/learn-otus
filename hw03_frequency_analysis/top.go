package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

type data struct {
	value string
	count int
}

func Top10(text string) []string {
	sl := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ',' || r == '\n' || r == '\t' || r == '.'
	})
	m := map[string]int{}
	for _, v := range sl {
		if v != "-" {
			m[strings.ToLower(v)]++
		}
	}
	fmt.Println(m)
	slData := make([]data, len(m))
	for i, v := range m {
		slData = append(slData, data{i, v})
	}

	sort.Slice(slData, func(i, j int) bool {
		if slData[i].count == slData[j].count {
			return strings.Compare(slData[i].value, slData[j].value) < 0
		} else {
			return slData[i].count > slData[j].count
		}
	})

	var index int = 10
	if index > len(slData) {
		index = len(slData)
	}
	result := make([]string, index)
	fmt.Println(index)
	for i := 0; i < index; i++ {
		result[i] = slData[i].value
	}
	fmt.Printf("%v", result)
	return result
}
