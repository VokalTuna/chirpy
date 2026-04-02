package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		secret         string
		validateSecret string
		duration       time.Duration
		expectError    bool
	}{
		{
			name:           "Happy path",
			userID:         uuid.New(),
			secret:         "mysecret",
			validateSecret: "mysecret",
			duration:       time.Hour,
			expectError:    false,
		},
		{
			name:           "Expired token",
			userID:         uuid.New(),
			secret:         "mysecret",
			validateSecret: "mysecret",
			duration:       -time.Second,
			expectError:    true,
		},
		{
			name:           "Wrong secret",
			userID:         uuid.New(),
			secret:         "mysecret",
			validateSecret: "wrongsecret",
			duration:       time.Hour,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.secret, tt.duration)
			if err != nil && !tt.expectError {
				t.Fatalf("MakeJWT() unexpected error: %v", err)
			}
			gotID, err := ValidateJWT(token, tt.validateSecret)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateJWT() error = %v, expectError %v", err, tt.expectError)
			}
			if !tt.expectError && gotID != tt.userID {
				t.Errorf("ValidateJWT() got %v, want %v", gotID, tt.userID)
			}
		})
	}
}
