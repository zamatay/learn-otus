package hw09structvalidator

import (
	"reflect"
	"strings"
)

type ValidatorItem struct {
	v          validStruct
	tag        string
	expression string
}

func (vi ValidatorItem) getValidStruct() validStruct {
	return vi.v
}

func (vi ValidatorItem) Validate(name string, reflectValue *reflect.Value) (*ValidationError, bool) {
	f, isOk := arrayFunc[vi.tag]
	if isOk {
		return f.fun(name, reflectValue, vi)
	}
	return nil, true
}

func (vi ValidatorItem) match(value string) (*ValidatorItem, bool) {
	idx := strings.IndexRune(value, ':')
	if idx == -1 {
		return nil, false
	}
	vi.tag, vi.expression = value[:idx], value[idx+1:]
	return &vi, true
}

type Validators struct {
	Items []ValidatorItemInt
}

func NewValidators() ValidatorsInt {
	i := make([]ValidatorItemInt, 0, 6)
	return &Validators{Items: i}
}

func (v *Validators) Add(item validStruct) bool {
	for _, vi := range v.Items {
		if vi.getValidStruct().name == item.name {
			return false
		}
	}
	v.Items = append(v.Items, ValidatorItem{v: item})
	return true
}

func (v *Validators) Match(data string) (*ValidatorItem, bool) {
	for _, vi := range v.Items {
		m, isOk := vi.match(data)
		if isOk && m.tag == vi.getValidStruct().name {
			m.v = vi.getValidStruct()
			return m, isOk
		}
	}
	return nil, false
}
