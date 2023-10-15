package hw09structvalidator

import "reflect"

type ValidatorItemInt interface {
	getValidStruct() validStruct
	match(v string) (*ValidatorItem, bool)
	Validate(name string, reflectValue *reflect.Value) (*ValidationError, bool)
}

type ValidatorsInt interface {
	Add(vs validStruct) bool
	Match(data string) (*ValidatorItem, bool)
}
