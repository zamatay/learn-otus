package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("dirNotFound", func(t *testing.T) {
		_, err := ReadDir("123")
		if !errors.Is(err, ErrDirIsNotLoad) {
			t.Fail()
		}
	})

	t.Run("dirFound", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		if err != nil {
			t.Fail()
		}
		if len(env) != len(envTest) {
			t.Fail()
		}
		for key, value := range env {
			if value.NeedRemove != envTest[key].NeedRemove || envTest[key].Value != value.Value {
				t.Fail()
			}
		}
		fmt.Printf("%v\n", env)
	})
}
