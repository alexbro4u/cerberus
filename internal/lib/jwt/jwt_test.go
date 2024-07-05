package jwt

import (
	"cerberus/internal/domain/models"
	"testing"
	"time"
)

func TestNewToken_WithValidUserAndApp_ReturnsToken(t *testing.T) {
	user := models.User{ID: 1234567890}
	app := models.App{ID: 12345, Secret: "testSecret"}
	duration := time.Hour

	token, err := NewToken(user, app, duration)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Errorf("Expected token, got an empty string")
	}
}

func TestNewToken_WithEmptySecret_ReturnsError(t *testing.T) {
	user := models.User{ID: 1234567890}
	app := models.App{ID: 12345, Secret: ""}
	duration := time.Hour

	_, err := NewToken(user, app, duration)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestNewToken_WithZeroDuration_ReturnsToken(t *testing.T) {
	user := models.User{ID: 1234567890}
	app := models.App{ID: 12345, Secret: "testSecret"}
	duration := time.Duration(0)

	token, err := NewToken(user, app, duration)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Errorf("Expected token, got an empty string")
	}
}
