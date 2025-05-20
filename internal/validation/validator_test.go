package validation

import (
	"strings"
	"testing"
)

func TestNameValidator(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid name",
			value:   "John Doe",
			wantErr: false,
		},
		{
			name:    "too long name",
			value:   strings.Repeat("a", GetConfig().MaxNameLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewNameValidator(tt.value)
			err := v.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NameValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if v.GetValue() != tt.value {
				t.Errorf("NameValidator.GetValue() = %v, want %v", v.GetValue(), tt.value)
			}
		})
	}
}

func TestAgeValidator(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid age",
			value:   25,
			wantErr: false,
		},
		{
			name:    "too young",
			value:   GetConfig().MinAge - 1,
			wantErr: true,
		},
		{
			name:    "too old",
			value:   GetConfig().MaxAge + 1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewAgeValidator(tt.value)
			err := v.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AgeValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if v.GetValue() != tt.value {
				t.Errorf("AgeValidator.GetValue() = %v, want %v", v.GetValue(), tt.value)
			}
		})
	}
}
