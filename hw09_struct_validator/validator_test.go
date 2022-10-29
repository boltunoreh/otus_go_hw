package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae-878b1e16adee",
				Name:   "Test",
				Age:    20,
				meta:   nil,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"89519111511", "89519111511"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae",
				Name:   "Test",
				Age:    20,
				meta:   nil,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"89519111511", "89519111511"},
			},
			expectedErr: ErrStrLen,
		},
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae-878b1e16adee",
				Name:   "Test",
				Age:    15,
				meta:   nil,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"89519111511", "89519111511"},
			},
			expectedErr: ErrIntMin,
		},
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae-878b1e16adee",
				Name:   "Test",
				Age:    60,
				meta:   nil,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"89519111511", "89519111511"},
			},
			expectedErr: ErrIntMax,
		},
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae-878b1e16adee",
				Name:   "Test",
				Age:    20,
				meta:   nil,
				Email:  "test-test.com",
				Role:   "admin",
				Phones: []string{"89519111511", "89519111511"},
			},
			expectedErr: ErrStrRegexp,
		},
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae-878b1e16adee",
				Name:   "Test",
				Age:    20,
				meta:   nil,
				Email:  "test@test.com",
				Role:   "guest",
				Phones: []string{"89519111511", "89519111511"},
			},
			expectedErr: ErrStrIn,
		},
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae-878b1e16adee",
				Name:   "Test",
				Age:    20,
				meta:   nil,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"9519111511", "89519111511"},
			},
			expectedErr: ErrStrLen,
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1234",
			},
			expectedErr: ErrStrLen,
		},
		{
			in: Token{
				Header:    []byte{1, 3, 4, 6},
				Payload:   []byte{1, 3, 4},
				Signature: []byte{1, 6},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 301,
				Body: "{\"redirect\": \"dzen.ru\"}",
			},
			expectedErr: ErrIntIn,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
