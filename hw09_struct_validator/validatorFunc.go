package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type execFunction func(fieldName string, reflectValue *reflect.Value, vi ValidatorItem) (*ValidationError, bool)

var (
	lenValidatorError    = errors.New("lenValidatorError")
	regexpValidatorError = errors.New("regexpValidatorError")
	minValidatorError    = errors.New("minValidatorError")
	maxValidatorError    = errors.New("maxValidatorError")
	inValidatorError     = errors.New("inValidatorError")
)
var validateLen = func(fieldName string, reflectValue *reflect.Value, vi ValidatorItem) (*ValidationError, bool) {
	if reflectValue.Kind() == 24 || reflectValue.Kind() == 23 {
		l, err := strconv.Atoi(vi.expression)
		if err == nil {
			values, isOk := reflectValue.Interface().([]string)
			if !isOk {
				values = append(values, reflectValue.String())
			}
			for _, value := range values {
				if len(value) > l {
					return &ValidationError{Field: fieldName, Err: vi.getValidStruct().err}, false
				}
			}
		}
	}
	return nil, true
}

var validateRegExp = func(fieldName string, reflectValue *reflect.Value, vi ValidatorItem) (*ValidationError, bool) {
	if reflectValue.Kind() == 24 || reflectValue.Kind() == 23 {
		rg, _ := regexp.Compile(vi.expression)
		value := reflectValue.String()
		if !rg.MatchString(value) {
			return &ValidationError{fieldName, vi.getValidStruct().err}, false
		}
	}
	return nil, true
}

var validateMinMax = func(fieldName string, reflectValue *reflect.Value, vi ValidatorItem) (*ValidationError, bool) {
	if (reflectValue.Kind() > 1 && reflectValue.Kind() < 15) || reflectValue.Kind() == 23 {
		l, err := strconv.Atoi(vi.expression)
		if err == nil {
			values, isOk := reflectValue.Interface().([]int64)
			if !isOk {
				values = append(values, reflectValue.Int())
			}
			for _, value := range values {
				if vi.tag == "min" {
					if value < int64(l) {
						return &ValidationError{fieldName, vi.v.err}, false
					}
				} else if vi.tag == "max" {
					if value > int64(l) {
						return &ValidationError{fieldName, vi.v.err}, false
					}
				}
			}
		}
	}
	return nil, true
}

var validateIn = func(fieldName string, reflectValue *reflect.Value, vi ValidatorItem) (*ValidationError, bool) {
	if (reflectValue.Kind() > 1 && reflectValue.Kind() < 15) || reflectValue.Kind() == 23 || reflectValue.Kind() == 24 {
		values := strings.Split(vi.expression, ",")
		if len(values) == 0 {
			return &ValidationError{fieldName, vi.getValidStruct().err}, false
		}

		vals, isOk := reflectValue.Interface().([]string)
		if !isOk {
			if reflectValue.Kind() == reflect.Int {
				vals = append(vals, strconv.FormatInt(reflectValue.Int(), 10))
			}
			if reflectValue.Kind() == reflect.String {
				vals = append(vals, reflectValue.String())
			}
		}

		for _, value := range vals {
			isOk := false
			for _, v := range values {
				isOk = isOk || (value == v)
				if isOk {
					break
				}
			}
			if !isOk {
				return &ValidationError{fieldName, vi.getValidStruct().err}, false
			}
		}
	}
	return nil, true
}

type validStruct struct {
	name string
	fun  execFunction
	err  error
}

var arrayFunc = map[string]validStruct{
	"len":    validStruct{"len", validateLen, lenValidatorError},
	"regexp": validStruct{"regexp", validateRegExp, regexpValidatorError},
	"min":    validStruct{name: "min", fun: validateMinMax, err: minValidatorError},
	"max":    validStruct{name: "max", fun: validateMinMax, err: maxValidatorError},
	"in":     validStruct{name: "in", fun: validateIn, err: inValidatorError},
}
