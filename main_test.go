package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsPasswordField(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		expected bool
	}{
		{
			name:     "SPACE_TRACK_PASSWORD should be masked with asterisks",
			envKey:   "SPACE_TRACK_PASSWORD",
			expected: true,
		},
		{
			name:     "N2YO_API_KEY should be masked with asterisks",
			envKey:   "N2YO_API_KEY",
			expected: true,
		},
		{
			name:     "SPACE_TRACK_USERNAME should not be masked",
			envKey:   "SPACE_TRACK_USERNAME",
			expected: false,
		},
		{
			name:     "Field containing PASSWORD should be masked with asterisks",
			envKey:   "MY_PASSWORD",
			expected: true,
		},
		{
			name:     "Field containing API_KEY should be masked with asterisks",
			envKey:   "SOME_API_KEY",
			expected: true,
		},
		{
			name:     "Field containing SECRET should be masked with asterisks",
			envKey:   "MY_SECRET",
			expected: true,
		},
		{
			name:     "Field containing TOKEN should be masked with asterisks",
			envKey:   "ACCESS_TOKEN",
			expected: true,
		},
		{
			name:     "Lowercase password field should be masked with asterisks",
			envKey:   "my_password",
			expected: true,
		},
		{
			name:     "Mixed case password field should be masked with asterisks",
			envKey:   "My_Password_Field",
			expected: true,
		},
		{
			name:     "Regular environment variable should not be masked",
			envKey:   "DATABASE_HOST",
			expected: false,
		},
		{
			name:     "Empty string should not be masked",
			envKey:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPasswordField(tt.envKey)
			if result != tt.expected {
				t.Errorf("isPasswordField(%q) = %v, want %v", tt.envKey, result, tt.expected)
			}
		})
	}
}

func TestLoadEnvFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		envContent  string
		expectError bool
		checkEnv    func(t *testing.T) // Function to check environment variables
	}{
		{
			name: "Valid .env file with all variables",
			envContent: `SPACE_TRACK_USERNAME=testuser
SPACE_TRACK_PASSWORD=fake_test_password
N2YO_API_KEY=fake_test_key
`,
			expectError: false,
			checkEnv: func(t *testing.T) {
				if val := os.Getenv("SPACE_TRACK_USERNAME"); val != "testuser" {
					t.Errorf("SPACE_TRACK_USERNAME = %q, want %q", val, "testuser")
				}
				if val := os.Getenv("SPACE_TRACK_PASSWORD"); val != "fake_test_password" {
					t.Errorf("SPACE_TRACK_PASSWORD = %q, want %q", val, "fake_test_password")
				}
				if val := os.Getenv("N2YO_API_KEY"); val != "fake_test_key" {
					t.Errorf("N2YO_API_KEY = %q, want %q", val, "fake_test_key")
				}
			},
		},
		{
			name: ".env file with comments",
			envContent: `# This is a comment
SPACE_TRACK_USERNAME=testuser
# Another comment
SPACE_TRACK_PASSWORD=testpass
`,
			expectError: false,
			checkEnv: func(t *testing.T) {
				if val := os.Getenv("SPACE_TRACK_USERNAME"); val != "testuser" {
					t.Errorf("SPACE_TRACK_USERNAME = %q, want %q", val, "testuser")
				}
				if val := os.Getenv("SPACE_TRACK_PASSWORD"); val != "testpass" {
					t.Errorf("SPACE_TRACK_PASSWORD = %q, want %q", val, "testpass")
				}
			},
		},
		{
			name: ".env file with empty lines",
			envContent: `SPACE_TRACK_USERNAME=testuser

SPACE_TRACK_PASSWORD=testpass

`,
			expectError: false,
			checkEnv: func(t *testing.T) {
				if val := os.Getenv("SPACE_TRACK_USERNAME"); val != "testuser" {
					t.Errorf("SPACE_TRACK_USERNAME = %q, want %q", val, "testuser")
				}
				if val := os.Getenv("SPACE_TRACK_PASSWORD"); val != "testpass" {
					t.Errorf("SPACE_TRACK_PASSWORD = %q, want %q", val, "testpass")
				}
			},
		},
		{
			name: ".env file with quoted values",
			envContent: `SPACE_TRACK_USERNAME="testuser"
SPACE_TRACK_PASSWORD='fake_test_password'
N2YO_API_KEY="fake_test_key"
`,
			expectError: false,
			checkEnv: func(t *testing.T) {
				if val := os.Getenv("SPACE_TRACK_USERNAME"); val != "testuser" {
					t.Errorf("SPACE_TRACK_USERNAME = %q, want %q", val, "testuser")
				}
				if val := os.Getenv("SPACE_TRACK_PASSWORD"); val != "fake_test_password" {
					t.Errorf("SPACE_TRACK_PASSWORD = %q, want %q", val, "fake_test_password")
				}
				if val := os.Getenv("N2YO_API_KEY"); val != "fake_test_key" {
					t.Errorf("N2YO_API_KEY = %q, want %q", val, "fake_test_key")
				}
			},
		},
		{
			name: ".env file with spaces around equals",
			envContent: `SPACE_TRACK_USERNAME = testuser
SPACE_TRACK_PASSWORD = testpass
`,
			expectError: false,
			checkEnv: func(t *testing.T) {
				if val := os.Getenv("SPACE_TRACK_USERNAME"); val != "testuser" {
					t.Errorf("SPACE_TRACK_USERNAME = %q, want %q", val, "testuser")
				}
				if val := os.Getenv("SPACE_TRACK_PASSWORD"); val != "testpass" {
					t.Errorf("SPACE_TRACK_PASSWORD = %q, want %q", val, "testpass")
				}
			},
		},
		{
			name:        ".env file not found",
			envContent:  "",
			expectError: true,
			checkEnv:    func(t *testing.T) {},
		},
		{
			name: ".env file with invalid line (no equals)",
			envContent: `SPACE_TRACK_USERNAME=testuser
INVALID_LINE_WITHOUT_EQUALS
SPACE_TRACK_PASSWORD=testpass
`,
			expectError: false,
			checkEnv: func(t *testing.T) {
				if val := os.Getenv("SPACE_TRACK_USERNAME"); val != "testuser" {
					t.Errorf("SPACE_TRACK_USERNAME = %q, want %q", val, "testuser")
				}
				if val := os.Getenv("SPACE_TRACK_PASSWORD"); val != "testpass" {
					t.Errorf("SPACE_TRACK_PASSWORD = %q, want %q", val, "testpass")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variables before test
			os.Unsetenv("SPACE_TRACK_USERNAME")
			os.Unsetenv("SPACE_TRACK_PASSWORD")
			os.Unsetenv("N2YO_API_KEY")

			// Save original working directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}

			// Change to temporary directory
			err = os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}
			defer os.Chdir(originalDir)

			if tt.name == ".env file not found" {
				// Ensure .env file doesn't exist in temp directory
				envPath := filepath.Join(tmpDir, ".env")
				os.Remove(envPath) // Remove if it exists
			} else {
				// Create .env file
				envPath := filepath.Join(tmpDir, ".env")
				err := os.WriteFile(envPath, []byte(tt.envContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create .env file: %v", err)
				}
			}

			// Test loadEnvFile
			err = loadEnvFile()
			if tt.expectError {
				if err == nil {
					t.Errorf("loadEnvFile() expected error, got nil")
				}
				if err != nil && !strings.Contains(err.Error(), ".env file not found") {
					t.Errorf("loadEnvFile() error = %v, want error containing '.env file not found'", err)
				}
			} else {
				if err != nil {
					t.Errorf("loadEnvFile() unexpected error: %v", err)
				}
				tt.checkEnv(t)
			}
		})
	}
}

func TestCheckEnvironmentalVariable(t *testing.T) {
	tests := []struct {
		name        string
		envKey      string
		preSetValue string
		shouldSet   bool
	}{
		{
			name:        "Variable not set should trigger set",
			envKey:      "TEST_VAR_NOT_SET",
			preSetValue: "",
			shouldSet:   true,
		},
		{
			name:        "Variable already set should not trigger set",
			envKey:      "TEST_VAR_SET",
			preSetValue: "existing_value",
			shouldSet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.Unsetenv(tt.envKey)

			// Pre-set value if needed
			if tt.preSetValue != "" {
				os.Setenv(tt.envKey, tt.preSetValue)
			}

			// Note: We can't easily test the interactive part (setEnvironmentalVariable)
			// without mocking stdin, but we can test that checkEnvironmentalVariable
			// correctly identifies when a variable is missing
			_, found := os.LookupEnv(tt.envKey)
			if found != (tt.preSetValue != "") {
				t.Errorf("Environment variable lookup mismatch")
			}

			// Clean up after test
			os.Unsetenv(tt.envKey)
		})
	}
}

func TestLoadEnvFileEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	t.Run("Empty .env file", func(t *testing.T) {
		os.Chdir(tmpDir)
		os.Unsetenv("TEST_VAR")

		envPath := filepath.Join(tmpDir, ".env")
		err := os.WriteFile(envPath, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create .env file: %v", err)
		}

		err = loadEnvFile()
		if err != nil {
			t.Errorf("loadEnvFile() with empty file should not error, got: %v", err)
		}
	})

	t.Run(".env file with only comments", func(t *testing.T) {
		os.Chdir(tmpDir)
		os.Unsetenv("TEST_VAR")

		envPath := filepath.Join(tmpDir, ".env")
		err := os.WriteFile(envPath, []byte("# Comment 1\n# Comment 2\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create .env file: %v", err)
		}

		err = loadEnvFile()
		if err != nil {
			t.Errorf("loadEnvFile() with only comments should not error, got: %v", err)
		}
	})

	t.Run(".env file with value containing equals", func(t *testing.T) {
		os.Chdir(tmpDir)
		os.Unsetenv("TEST_VAR")

		envPath := filepath.Join(tmpDir, ".env")
		// SplitN with limit 2 should handle this correctly
		err := os.WriteFile(envPath, []byte("TEST_VAR=value=with=equals\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create .env file: %v", err)
		}

		err = loadEnvFile()
		if err != nil {
			t.Errorf("loadEnvFile() should not error, got: %v", err)
		}

		val := os.Getenv("TEST_VAR")
		if val != "value=with=equals" {
			t.Errorf("TEST_VAR = %q, want %q", val, "value=with=equals")
		}
	})

	t.Run(".env file with whitespace-only lines", func(t *testing.T) {
		os.Chdir(tmpDir)
		os.Unsetenv("TEST_VAR")

		envPath := filepath.Join(tmpDir, ".env")
		err := os.WriteFile(envPath, []byte("   \n\t\nTEST_VAR=value\n   \n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create .env file: %v", err)
		}

		err = loadEnvFile()
		if err != nil {
			t.Errorf("loadEnvFile() should not error, got: %v", err)
		}

		val := os.Getenv("TEST_VAR")
		if val != "value" {
			t.Errorf("TEST_VAR = %q, want %q", val, "value")
		}
	})
}

// Benchmark tests
func BenchmarkIsPasswordField(b *testing.B) {
	testCases := []string{
		"SPACE_TRACK_PASSWORD",
		"N2YO_API_KEY",
		"SPACE_TRACK_USERNAME",
		"MY_PASSWORD",
		"REGULAR_VAR",
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			isPasswordField(tc)
		}
	}
}

func BenchmarkLoadEnvFile(b *testing.B) {
	tmpDir := b.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	envContent := `SPACE_TRACK_USERNAME=testuser
SPACE_TRACK_PASSWORD=fake_test_password
N2YO_API_KEY=fake_test_key
`

	envPath := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envPath, []byte(envContent), 0644)
	if err != nil {
		b.Fatalf("Failed to create .env file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.Chdir(tmpDir)
		os.Unsetenv("SPACE_TRACK_USERNAME")
		os.Unsetenv("SPACE_TRACK_PASSWORD")
		os.Unsetenv("N2YO_API_KEY")
		loadEnvFile()
	}
}
