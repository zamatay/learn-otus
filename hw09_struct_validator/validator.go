package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
)

var ValidateType = []reflect.Kind{
	reflect.String,
	reflect.Int,
	reflect.Slice,
}

var rgLen *regexp.Regexp
var rgRegexp *regexp.Regexp
var rgIn *regexp.Regexp
var rgMinMax *regexp.Regexp

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func init() {
	rgLen = regexp.MustCompile(`len:(\d+)`)
	rgRegexp = regexp.MustCompile(`regexp:*`)
	rgIn = regexp.MustCompile(`in:*`)
	rgMinMax = regexp.MustCompile(`(min|max):(\d+)`)
}

func Validate(v interface{}) error {

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return errors.New("NoStruct")
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f.Type.Kind())

		if slices.Contains(ValidateType, f.Type.Kind()) {
			validateStr := f.Tag.Get("validate")
			if validateStr != "" {
				validateField(validateStr)
			}
		}
	}
	return nil
}

func validateField(str string) {

}
