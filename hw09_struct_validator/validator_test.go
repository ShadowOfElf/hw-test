package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID       string `json:"id" validate:"len:36"`
		Name     string
		Age      int             `validate:"min:18|max:50"`
		Email    string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role     UserRole        `validate:"in:admin,stuff"`
		Phones   []string        `validate:"len:11"`
		Response Response        `validate:"struct:true"`
		meta     json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	MapErr1 struct {
		Text string `validate:"len5"`
	}

	MapErr2 struct {
		Text string `validate:"len:5:6"`
	}

	MapErr3 struct {
		Text string `validate:"max:6"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:       "10",
				Name:     "Test",
				Age:      50,
				Email:    "a@a.ru",
				Role:     "admin",
				Response: Response{Code: 200, Body: "test"},
				Phones:   []string{"12345678910", "911", "0000000"},
			},
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in:          App{Version: "vers"},
			expectedErr: nil,
		},
		{
			in:          22,
			expectedErr: ErrWrongArgument,
		},
		{
			in: User{
				ID:       "10",
				Name:     "Test",
				Age:      51,
				Email:    "a@a.ru",
				Role:     "admin",
				Response: Response{Code: 200, Body: "test"},
				Phones:   []string{"12345678910", "911", "0000000"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrValidateMax},
			},
		},
		{
			in: MapErr1{Text: "text_test"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Text", Err: ErrWrongFormat},
			},
		},
		{
			in: MapErr2{Text: "text_test"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Text", Err: ErrWrongFormat},
			},
		},
		{
			in: MapErr3{Text: "text_test"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Text", Err: ErrWrongParamType},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			var validationErrs ValidationErrors
			var expected ValidationErrors
			if errors.As(err, &validationErrs) {
				errors.As(tt.expectedErr, &expected)
				require.ErrorIs(t, expected[0].Err, validationErrs[0].Err)
			} else {
				require.ErrorIs(t, tt.expectedErr, err)
			}
			_ = tt
		})
	}
}
