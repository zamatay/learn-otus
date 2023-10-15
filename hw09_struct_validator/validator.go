package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

var (
	ValidateType = []reflect.Kind{
		reflect.String,
		reflect.Int,
		reflect.Slice,
	}
	Validator      = NewValidators().(*Validators)
	NoStructErrors = errors.New("NoStruct")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errorText strings.Builder
	for idx, field := range v {
		if idx != 0 {
			errorText.WriteRune('\n')
		}
		errorText.WriteString(fmt.Sprintf("%s: %s", field.Field, field.Err.Error()))
	}
	return errorText.String()
}

func init() {
	Validator.Add(arrayFunc["len"])
	Validator.Add(arrayFunc["regexp"])
	Validator.Add(arrayFunc["min"])
	Validator.Add(arrayFunc["max"])
	Validator.Add(arrayFunc["in"])
}

func Validate(v interface{}) error {
	tt := reflect.TypeOf(v)
	if tt.Kind() != reflect.Struct {
		return NoStructErrors
	}
	tv := reflect.ValueOf(v)
	return handleValidate(tt, tv)
}

func handleValidate(tt reflect.Type, tv reflect.Value) error {
	validErrors := make(ValidationErrors, 0)
	for i := 0; i < tt.NumField(); i++ {
		ft := tt.Field(i)
		if slices.Contains(ValidateType, ft.Type.Kind()) {
			validateStr := ft.Tag.Get("validate")
			if validateStr != "" {
				fv := tv.Field(i)
				vf, isValid := validateField(ft.Name, &fv, validateStr)
				if !isValid {
					for _, v := range vf {
						validErrors = append(validErrors, v)
					}
				}
			}
		}
	}
	if len(validErrors) == 0 {
		return nil
	}
	return validErrors
}

func validateField(fieldName string, f *reflect.Value, tagExpression string) ([]ValidationError, bool) {
	expressions := strings.Split(tagExpression, "|")
	if len(expressions) == 0 {
		expressions = append(expressions, tagExpression)
	}
	er := make([]ValidationError, 0, 1)
	for _, ex := range expressions {
		m, isOk := Validator.Match(ex)
		if isOk {
			err, isOk := m.Validate(fieldName, f)
			if !isOk {
				er = append(er, *err)
			}
		}
	}
	if len(er) > 0 {
		return er, false
	}
	return nil, true
}
