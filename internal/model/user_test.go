package model

import (
	"testing"
)

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
