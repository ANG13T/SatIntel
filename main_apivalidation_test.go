package main

// NOTE: This test file contains fake test credentials (e.g., "fake_test_password", "FAKE-TEST-KEY")
// These are NOT real secrets and are used only for testing purposes.

import (
	"os"
	"strings"
	"testing"
)

func TestValidateAPIKeyFormat(t *testing.T) {
	tests := []struct {
		name        string
		envKey      string
		value       string
		expectError bool
		errorMsg    string
	}{
		// Valid cases
		{
			name:        "Valid Space-Track username",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "testuser123",
			expectError: false,
		},
		{
			name:        "Valid Space-Track password",
			envKey:      "SPACE_TRACK_PASSWORD",
			value:       "fake_test_password_123", // Test value only
			expectError: false,
		},
		{
			name:        "Valid N2YO API key",
			envKey:      "N2YO_API_KEY",
			value:       "FAKE-TEST-KEY-12345-67890", // Test value only
			expectError: false,
		},
		{
			name:        "Valid N2YO API key with underscores",
			envKey:      "N2YO_API_KEY",
			value:       "test_key_123",
			expectError: false,
		},
		{
			name:        "Valid long credentials",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       strings.Repeat("a", 100),
			expectError: false,
		},
		// Empty values
		{
			name:        "Empty username",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "Whitespace only username",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "   ",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "Empty password",
			envKey:      "SPACE_TRACK_PASSWORD",
			value:       "",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "Empty API key",
			envKey:      "N2YO_API_KEY",
			value:       "",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		// Too short
		{
			name:        "Username too short",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "ab",
			expectError: true,
			errorMsg:    "too short",
		},
		{
			name:        "Password too short",
			envKey:      "SPACE_TRACK_PASSWORD",
			value:       "12",
			expectError: true,
			errorMsg:    "too short",
		},
		{
			name:        "API key too short",
			envKey:      "N2YO_API_KEY",
			value:       "ab",
			expectError: true,
			errorMsg:    "too short",
		},
		// Too long
		{
			name:        "Username too long",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       strings.Repeat("a", 513),
			expectError: true,
			errorMsg:    "too long",
		},
		{
			name:        "Password too long",
			envKey:      "SPACE_TRACK_PASSWORD",
			value:       strings.Repeat("a", 513),
			expectError: true,
			errorMsg:    "too long",
		},
		{
			name:        "API key too long",
			envKey:      "N2YO_API_KEY",
			value:       strings.Repeat("a", 513),
			expectError: true,
			errorMsg:    "too long",
		},
		// Invalid characters
		{
			name:        "Username with newline",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "user\nname",
			expectError: true,
			errorMsg:    "invalid characters",
		},
		{
			name:        "Username with carriage return",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "user\rname",
			expectError: true,
			errorMsg:    "invalid characters",
		},
		{
			name:        "Username with tab",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "user\tname",
			expectError: true,
			errorMsg:    "invalid characters",
		},
		{
			name:        "API key with space",
			envKey:      "N2YO_API_KEY",
			value:       "key with spaces",
			expectError: true,
			errorMsg:    "invalid whitespace",
		},
		{
			name:        "API key with newline",
			envKey:      "N2YO_API_KEY",
			value:       "key\nwith\nnewlines",
			expectError: true,
			errorMsg:    "invalid whitespace",
		},
		{
			name:        "API key with only special chars",
			envKey:      "N2YO_API_KEY",
			value:       "---",
			expectError: true,
			errorMsg:    "must contain at least one alphanumeric",
		},
		// Edge cases
		{
			name:        "Minimum length username",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       "abc",
			expectError: false,
		},
		{
			name:        "Maximum length username",
			envKey:      "SPACE_TRACK_USERNAME",
			value:       strings.Repeat("a", 512),
			expectError: false,
		},
		{
			name:        "API key with mixed case",
			envKey:      "N2YO_API_KEY",
			value:       "FAKE-TEST-KEY-ABC-XYZ", // Test value only
			expectError: false,
		},
		{
			name:        "API key with numbers only",
			envKey:      "N2YO_API_KEY",
			value:       "123456789",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAPIKeyFormat(tt.envKey, tt.value)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateAPIKeyFormat(%q, %q) expected error, got nil", tt.envKey, tt.value)
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateAPIKeyFormat(%q, %q) error = %v, want error containing %q", tt.envKey, tt.value, err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateAPIKeyFormat(%q, %q) unexpected error: %v", tt.envKey, tt.value, err)
				}
			}
		})
	}
}

func TestValidateAPIKeyFormatUnknownKey(t *testing.T) {
	// Test with unknown environment variable key
	err := validateAPIKeyFormat("UNKNOWN_KEY", "testvalue")
	if err != nil {
		t.Errorf("validateAPIKeyFormat with unknown key should only check basic format, got error: %v", err)
	}
}

func TestValidateCredentialsFormatOnly(t *testing.T) {
	// Test format validation without actual API calls
	// This requires mocking or skipping connection tests

	// Save original values
	origUsername := os.Getenv("SPACE_TRACK_USERNAME")
	origPassword := os.Getenv("SPACE_TRACK_PASSWORD")
	origAPIKey := os.Getenv("N2YO_API_KEY")

	// Set test values
	os.Setenv("SPACE_TRACK_USERNAME", "testuser")
		os.Setenv("SPACE_TRACK_PASSWORD", "fake_test_password") // Test value only
		os.Setenv("N2YO_API_KEY", "fake_test_key_123") // Test value only

	// Test format validation (we'll skip connection tests in unit tests)
	username := os.Getenv("SPACE_TRACK_USERNAME")
	password := os.Getenv("SPACE_TRACK_PASSWORD")
	apiKey := os.Getenv("N2YO_API_KEY")

	if err := validateAPIKeyFormat("SPACE_TRACK_USERNAME", username); err != nil {
		t.Errorf("Username format validation failed: %v", err)
	}

	if err := validateAPIKeyFormat("SPACE_TRACK_PASSWORD", password); err != nil {
		t.Errorf("Password format validation failed: %v", err)
	}

	if err := validateAPIKeyFormat("N2YO_API_KEY", apiKey); err != nil {
		t.Errorf("API key format validation failed: %v", err)
	}

	// Restore original values
	if origUsername != "" {
		os.Setenv("SPACE_TRACK_USERNAME", origUsername)
	} else {
		os.Unsetenv("SPACE_TRACK_USERNAME")
	}
	if origPassword != "" {
		os.Setenv("SPACE_TRACK_PASSWORD", origPassword)
	} else {
		os.Unsetenv("SPACE_TRACK_PASSWORD")
	}
	if origAPIKey != "" {
		os.Setenv("N2YO_API_KEY", origAPIKey)
	} else {
		os.Unsetenv("N2YO_API_KEY")
	}
}

// TestValidateCredentialsFormatOnlyWithInvalidInput tests format validation with invalid inputs
func TestValidateCredentialsFormatOnlyWithInvalidInput(t *testing.T) {
	// Save original values
	origUsername := os.Getenv("SPACE_TRACK_USERNAME")
	origPassword := os.Getenv("SPACE_TRACK_PASSWORD")
	origAPIKey := os.Getenv("N2YO_API_KEY")

	defer func() {
		if origUsername != "" {
			os.Setenv("SPACE_TRACK_USERNAME", origUsername)
		} else {
			os.Unsetenv("SPACE_TRACK_USERNAME")
		}
		if origPassword != "" {
			os.Setenv("SPACE_TRACK_PASSWORD", origPassword)
		} else {
			os.Unsetenv("SPACE_TRACK_PASSWORD")
		}
		if origAPIKey != "" {
			os.Setenv("N2YO_API_KEY", origAPIKey)
		} else {
			os.Unsetenv("N2YO_API_KEY")
		}
	}()

	// Test with invalid format credentials
	os.Setenv("SPACE_TRACK_USERNAME", "ab") // Too short
		os.Setenv("SPACE_TRACK_PASSWORD", "fake_test_password") // Test value only
		os.Setenv("N2YO_API_KEY", "fake_test_key_123") // Test value only

	username := os.Getenv("SPACE_TRACK_USERNAME")
	if err := validateAPIKeyFormat("SPACE_TRACK_USERNAME", username); err == nil {
		t.Error("Expected error for too short username, got nil")
	}

	// Test with valid format
	os.Setenv("SPACE_TRACK_USERNAME", "testuser")
	if err := validateAPIKeyFormat("SPACE_TRACK_USERNAME", "testuser"); err != nil {
		t.Errorf("Expected no error for valid username, got: %v", err)
	}
}

// TestValidateAPIKeyFormatPasswordField tests password field validation
func TestValidateAPIKeyFormatPasswordField(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "Valid password",
			value:       "fake_test_pass_123", // Test value only
			expectError: false,
		},
		{
			name:        "Password with special characters",
			value:       "Pass@123!",
			expectError: false, // Passwords can have special chars
		},
		{
			name:        "Password too short",
			value:       "ab",
			expectError: true,
		},
		{
			name:        "Password too long",
			value:       strings.Repeat("a", 513),
			expectError: true,
		},
		{
			name:        "Empty password",
			value:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAPIKeyFormat("SPACE_TRACK_PASSWORD", tt.value)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %q, got nil", tt.value)
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %q: %v", tt.value, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidateAPIKeyFormat(b *testing.B) {
	testCases := []struct {
		envKey string
		value  string
	}{
		{"SPACE_TRACK_USERNAME", "testuser123"},
		{"SPACE_TRACK_PASSWORD", "fake_test_password_123"}, // Test value only
		{"N2YO_API_KEY", "FAKE-TEST-KEY-12345"}, // Test value only
		{"SPACE_TRACK_USERNAME", strings.Repeat("a", 100)},
		{"N2YO_API_KEY", "fake_test_key_123"}, // Test value only
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			validateAPIKeyFormat(tc.envKey, tc.value)
		}
	}
}
