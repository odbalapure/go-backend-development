package val

import (
	"testing"
)

func TestValidateString(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		minLength int
		maxLength int
		wantErr   bool
	}{
		{"valid length", "hello", 3, 10, false},
		{"too short", "hi", 3, 10, true},
		{"too long", "this is a very long string", 3, 10, true},
		{"exact min length", "abc", 3, 10, false},
		{"exact max length", "abcdefghij", 3, 10, false},
		{"empty string", "", 0, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateString(tt.value, tt.minLength, tt.maxLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid username", "john_doe123", false},
		{"valid username with underscore", "user_name", false},
		{"valid username with numbers", "user123", false},
		{"too short", "ab", true},
		{"too long", "a" + string(make([]byte, 100)), true},
		{"contains uppercase", "JohnDoe", true},
		{"contains special chars", "john@doe", true},
		{"contains spaces", "john doe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFullName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid full name", "John Doe", false},
		{"valid single name", "John", false},
		{"valid name with multiple spaces", "John Michael Doe", false},
		{"too short", "Jo", true},
		{"too long", "a" + string(make([]byte, 100)), true},
		{"contains numbers", "John123", true},
		{"contains special chars", "John@Doe", true},
		{"contains underscore", "John_Doe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFullName(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFullName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid password", "password123", false},
		{"valid password with special chars", "p@ssw0rd!", false},
		{"too short", "12345", true},
		{"exact min length", "123456", false},
		{"exact max length", string(make([]byte, 100)), false},
		{"too long", string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid email", "john@example.com", false},
		{"valid email with subdomain", "john@sub.example.com", false},
		{"valid email with plus", "john+tag@example.com", false},
		{"too short", "a@b", false}, // 3 chars is valid
		{"too long", "a" + string(make([]byte, 200)) + "@example.com", true},
		{"invalid email format", "invalid-email", true},
		{"missing @", "johnexample.com", true},
		{"missing domain", "john@", true},
		{"missing local part", "@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmailId(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{"valid positive id", 1, false},
		{"valid large id", 999999, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"large negative", -999999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailId(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmailId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSecretCode(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid secret code", string(make([]byte, 32)), false},
		{"valid long secret code", string(make([]byte, 64)), false},
		{"exact min length", string(make([]byte, 32)), false},
		{"exact max length", string(make([]byte, 128)), false},
		{"too short", string(make([]byte, 31)), true},
		{"too long", string(make([]byte, 129)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSecretCode(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSecretCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
