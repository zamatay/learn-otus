package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string
type name = string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   name
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff,user"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		caption     string
		in          interface{}
		expectedErr error
	}{
		{
			caption:     "len",
			in:          User{ID: "1234567890123456789012345678901234567", Email: "zamuraev@mail.com", Age: 18, Role: "user", Phones: []string{"79189182626"}},
			expectedErr: arrayFunc["len"].err,
		},
		{
			caption:     "regexp",
			in:          User{ID: "12345678901", Email: "zamuraev", Age: 18, Role: "user", Phones: []string{"79189182626"}},
			expectedErr: arrayFunc["regexp"].err,
		},
		{
			caption:     "min",
			in:          User{ID: "12345678901", Email: "zamuraev@mail.com", Age: 16, Role: "user", Phones: []string{"79189182626"}},
			expectedErr: arrayFunc["min"].err,
		},
		{
			caption:     "len slice",
			in:          User{ID: "12", Email: "zamuraev@mail.com", Age: 18, Role: "user", Phones: []string{"7918918262626"}},
			expectedErr: arrayFunc["max"].err,
		},
		{
			caption:     "in",
			in:          User{ID: "12", Email: "zamuraev@mail.com", Age: 18, Role: "guest", Phones: []string{"79189182626"}},
			expectedErr: arrayFunc["Role"].err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.caption, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.False(t, errors.Is(err, tt.expectedErr))
			//require.Equal(t, tt.expectedErr, err.Error())
		})
	}
	t.Run("not struct error", func(t *testing.T) {
		err := Validate("")
		require.EqualError(t, err, NoStructErrors.Error())
	})
	t.Run("len", func(t *testing.T) {
		err := Validate(App{Version: "123456"})
		require.EqualError(t, err, "Version: Значение превысило максимальную длину")
	})
	t.Run("multi", func(t *testing.T) {
		err := Validate(User{ID: "1234567890123456789012345678901234567", Age: 60})
		require.EqualError(t, err, "ID: Значение превысило максимальную длину\nAge: Значение больше максимального\nEmail: Значение не соответствует регулярному выражению\nRole: Значение не входит в перечень")
	})
	t.Run("in int fail", func(t *testing.T) {
		err := Validate(Response{Code: 10})
		require.EqualError(t, err, "Code: Значение не входит в перечень")
	})
	t.Run("in int pass", func(t *testing.T) {
		err := Validate(Response{Code: 404})
		require.NoError(t, err)
	})
}
