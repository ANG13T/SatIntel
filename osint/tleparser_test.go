package osint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid relative path",
			path:        "test.tle",
			expectError: false,
		},
		{
			name:        "Valid path with subdirectory",
			path:        "data/test.tle",
			expectError: false,
		},
		{
			name:        "Empty path",
			path:        "",
			expectError: true,
			errorMsg:    "file path cannot be empty",
		},
		{
			name:        "Whitespace only path",
			path:        "   ",
			expectError: true,
			errorMsg:    "file path cannot be empty",
		},
		{
			name:        "Path with null byte",
			path:        "test\x00.tle",
			expectError: true,
			errorMsg:    "file path contains invalid characters",
		},
		{
			name:        "Directory traversal with ../",
			path:        "../test.tle",
			expectError: true,
			errorMsg:    "directory traversal detected",
		},
		{
			name:        "Directory traversal with ..\\",
			path:        "..\\test.tle",
			expectError: true,
			errorMsg:    "directory traversal detected",
		},
		{
			name:        "Directory traversal in middle",
			path:        "data/../test.tle",
			expectError: true,
			errorMsg:    "directory traversal detected",
		},
		{
			name:        "Multiple directory traversals",
			path:        "../../../etc/passwd",
			expectError: true,
			errorMsg:    "directory traversal detected",
		},
		{
			name:        "Path with double slashes",
			path:        "data//test.tle",
			expectError: true,
			errorMsg:    "directory traversal detected",
		},
		{
			name:        "Path with double backslashes",
			path:        "data\\\\test.tle",
			expectError: true,
			errorMsg:    "directory traversal detected",
		},
		{
			name:        "Path too long",
			path:        strings.Repeat("a", 4097),
			expectError: true,
			errorMsg:    "file path is too long",
		},
		{
			name:        "Valid absolute path (Unix)",
			path:        "/tmp/test.tle",
			expectError: false,
		},
		{
			name:        "Valid absolute path (Windows)",
			path:        "C:\\Users\\test.tle",
			expectError: false,
		},
		{
			name:        "Path with spaces",
			path:        "my file.tle",
			expectError: false,
		},
		{
			name:        "Path with special characters",
			path:        "test-file_123.tle",
			expectError: false,
		},
		{
			name:        "Encoded directory traversal",
			path:        "%2e%2e%2f",
			expectError: false, // URL encoding not decoded, but still dangerous
			// Note: In production, you might want to decode and check
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateFilePath(%q) expected error, got nil", tt.path)
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateFilePath(%q) error = %v, want error containing %q", tt.path, err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateFilePath(%q) unexpected error: %v", tt.path, err)
				}
			}
		})
	}
}

func TestValidateFilePathWithRealFile(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tle")
	
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with valid file path
	err = validateFilePath(testFile)
	if err != nil {
		t.Errorf("validateFilePath(%q) with valid file got error: %v", testFile, err)
	}

	// Test with relative path to the temp file
	relPath, err := filepath.Rel(".", testFile)
	if err == nil {
		err = validateFilePath(relPath)
		if err != nil {
			t.Logf("Note: Relative path validation may fail if path contains '..': %v", err)
		}
	}
}

func TestValidateFilePathEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "Single dot",
			path:        ".",
			expectError: false, // filepath.Clean will handle this
		},
		{
			name:        "Current directory reference",
			path:        "./test.tle",
			expectError: false, // filepath.Clean will normalize this
		},
		{
			name:        "Path starting with dot",
			path:        ".test.tle",
			expectError: false, // Hidden file, but valid
		},
		{
			name:        "Path with many slashes",
			path:        "a/b/c/d/e/f/test.tle",
			expectError: false,
		},
		{
			name:        "Path at max length",
			path:        strings.Repeat("a", 4096),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path)
			if tt.expectError && err == nil {
				t.Errorf("validateFilePath(%q) expected error, got nil", tt.path)
			} else if !tt.expectError && err != nil {
				t.Errorf("validateFilePath(%q) unexpected error: %v", tt.path, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidateFilePath(b *testing.B) {
	testPaths := []string{
		"test.tle",
		"data/test.tle",
		"../test.tle",
		strings.Repeat("a", 100),
		"test\x00.tle",
	}

	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			validateFilePath(path)
		}
	}
}

