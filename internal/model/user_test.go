package model

import (
	"testing"

	"github.com/fatih/color"
)

func init() {
	color.NoColor = false
}

func TestValidadeEmail(t *testing.T) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	testCases := []struct {
		name      string
		email     string
		expectErr bool
		errText   string
	}{
		{
			name:      "Invalid Email",
			email:     "@.com",
			expectErr: true,
			errText:   "email deve ter o formato 'exemplo@dominio.com'",
		},
		{
			name:      "Valid Email with & char",
			email:     "testemail&@text.com",
			expectErr: false,
			errText:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateEmail(tc.email)
			if err != nil && !tc.expectErr {
				t.Errorf(red("unexpected error %v"), err)
			}
			if err == nil && tc.expectErr {
				t.Errorf(red("unexpected error %v"), err)
			}
			if !tc.expectErr && err == nil {
				t.Logf(green("PASS: %s"), tc.name)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testeCases := []struct {
		name      string
		user      User
		expectErr bool
		errText   string
	}{
		{
			name: "Valid user",
			user: User{
				Name:     "valid",
				Phone:    "62995032333",
				Password: "valid",
			},
			expectErr: false,
			errText:   "",
		},
		{
			name: "Empty name ",
			user: User{
				Name:     "",
				Phone:    "62995032333",
				Password: "validPassword",
			},
			expectErr: true,
			errText:   "Nome não pode ser vazio",
		},
	}

	for _, tc := range testeCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.Validate()
			if err != nil && !tc.expectErr {
				t.Errorf("unexpected error %v", err)
			}
			if err == nil && tc.expectErr {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}
