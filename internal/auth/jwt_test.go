package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	testUUID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("Failed to parse UUID: %v", err)
	}

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   int64
		expectedId  uuid.UUID
		expectedErr bool
	}{
		{
			name:        "valid token",
			userID:      testUUID,
			tokenSecret: "mysecret",
			expiresIn:   3600,
			expectedId:  testUUID,
			expectedErr: false,
		},
		{
			name:        "expired token",
			userID:      testUUID,
			tokenSecret: "mysecret",
			expiresIn:   -3600,
			expectedId:  uuid.Nil,
			expectedErr: true,
		},
		{
			name:        "invalid secret",
			userID:      testUUID,
			tokenSecret: "wrongsecret",
			expiresIn:   3600,
			expectedId:  uuid.Nil,
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.tokenSecret, time.Duration(tt.expiresIn)*time.Second)
			if err != nil {
				t.Fatalf("MakeJWT() error = %v", err)
			}

			userID, err := ValidateJWT(token, "mysecret")
			if (err != nil) != tt.expectedErr {
				t.Errorf("ValidateJWT() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if userID != tt.expectedId {
				t.Errorf("ValidateJWT() got = %v, want %v", userID, tt.expectedId)
			}
		})
	}
}
