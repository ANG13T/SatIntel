/*
Copyright © 2023 Angelina Tsuboi angelinatsuboi@proton.me
*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/ANG13T/SatIntel/cli"
	"golang.org/x/term"
)

// loadEnvFile reads environment variables from a .env file in the current directory.
// It skips empty lines and comments, and handles quoted values.
func loadEnvFile() error {
	envPath := ".env"

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return fmt.Errorf(".env file not found")
	}

	file, err := os.Open(envPath)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"'")
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

// isPasswordField determines if the environment variable should have masked input.
func isPasswordField(envKey string) bool {
	passwordFields := []string{
		"SPACE_TRACK_PASSWORD",
		"N2YO_API_KEY",
		"PASSWORD",
		"API_KEY",
		"SECRET",
		"TOKEN",
	}
	envKeyUpper := strings.ToUpper(envKey)
	for _, field := range passwordFields {
		if strings.Contains(envKeyUpper, field) {
			return true
		}
	}
	return false
}

// readPassword reads a password from stdin and displays asterisks for each character typed.
func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		// Fallback to regular input if not a terminal
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(input), nil
	}

	// Save current terminal state and set to raw mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("failed to set raw terminal: %w", err)
	}
	defer term.Restore(fd, oldState)

	var password []byte
	var input [1]byte

	for {
		n, err := os.Stdin.Read(input[:])
		if err != nil || n == 0 {
			break
		}

		char := input[0]

		// Handle Enter key (carriage return or newline)
		if char == '\r' || char == '\n' {
			fmt.Println()
			break
		}

		// Handle Backspace/Delete (127 = DEL, 8 = BS)
		if char == 127 || char == 8 {
			if len(password) > 0 {
				password = password[:len(password)-1]
				// Move cursor back, print space, move cursor back again
				fmt.Print("\b \b")
			}
			continue
		}

		// Handle Ctrl+C
		if char == 3 {
			fmt.Println()
			os.Exit(1)
		}

		// Skip control characters except those we handle
		if char < 32 {
			continue
		}

		// Add character to password and print asterisk
		password = append(password, char)
		fmt.Print("*")
	}

	return string(password), nil
}

// setEnvironmentalVariable prompts the user to enter a value for the given environment variable.
// It reads from stdin and sets the environment variable with the provided value.
// Password fields display asterisks (*) for each character typed for security.
func setEnvironmentalVariable(envKey string) string {
	var input string
	var err error

	for {
		fmt.Printf("%s: ", envKey)

		if isPasswordField(envKey) {
			input, err = readPassword()
			if err != nil {
				fmt.Println("Error reading password:", err)
				os.Exit(1)
			}
		} else {
			reader := bufio.NewReader(os.Stdin)
			input, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				os.Exit(1)
			}
			input = strings.TrimSpace(input)
		}

		input = strings.TrimSpace(input)

		// Validate format
		if err := validateAPIKeyFormat(envKey, input); err != nil {
			fmt.Printf("  [!] Validation error: %v\n", err)
			fmt.Println("Please enter a valid value:")
			continue
		}

		break
	}

	if err := os.Setenv(envKey, input); err != nil {
		fmt.Printf("Error setting environment variable %s: %v\n", envKey, err)
		os.Exit(1)
	}

	return input
}

// checkEnvironmentalVariable verifies if an environment variable exists.
// If it doesn't exist, it prompts the user to provide a value.
func checkEnvironmentalVariable(envKey string) {
	_, found := os.LookupEnv(envKey)
	if !found {
		setEnvironmentalVariable(envKey)
	}
}

// validateAPIKeyFormat validates the format of API keys and credentials.
// Returns an error if the format is invalid.
func validateAPIKeyFormat(envKey, value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return fmt.Errorf("%s cannot be empty", envKey)
	}

	// Check for reasonable length limits
	if len(value) < 3 {
		return fmt.Errorf("%s is too short (minimum 3 characters)", envKey)
	}

	if len(value) > 512 {
		return fmt.Errorf("%s is too long (maximum 512 characters)", envKey)
	}

	// Space-Track username validation
	if envKey == "SPACE_TRACK_USERNAME" {
		// Username should not contain certain special characters
		if strings.ContainsAny(value, "\n\r\t") {
			return fmt.Errorf("username contains invalid characters")
		}
	}

	// N2YO API key validation (typically alphanumeric, may contain hyphens)
	if envKey == "N2YO_API_KEY" {
		// N2YO API keys are typically alphanumeric with possible hyphens/underscores
		// Allow alphanumeric, hyphens, underscores
		hasAlphanumeric := false
		for _, r := range value {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				hasAlphanumeric = true
			} else if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return fmt.Errorf("API key contains invalid whitespace characters")
			} else if r != '-' && r != '_' {
				// Allow hyphens and underscores, but warn about other special chars
				// (We'll be lenient here as API key formats can vary)
			}
		}
		if !hasAlphanumeric {
			return fmt.Errorf("API key must contain at least one alphanumeric character")
		}
	}

	return nil
}

// validateCredentials validates all API credentials and tests connections.
func validateCredentials() error {
	username := os.Getenv("SPACE_TRACK_USERNAME")
	password := os.Getenv("SPACE_TRACK_PASSWORD")
	apiKey := os.Getenv("N2YO_API_KEY")

	// Validate format
	if err := validateAPIKeyFormat("SPACE_TRACK_USERNAME", username); err != nil {
		return fmt.Errorf("Space-Track username validation failed: %w", err)
	}

	if err := validateAPIKeyFormat("SPACE_TRACK_PASSWORD", password); err != nil {
		return fmt.Errorf("Space-Track password validation failed: %w", err)
	}

	if err := validateAPIKeyFormat("N2YO_API_KEY", apiKey); err != nil {
		return fmt.Errorf("N2YO API key validation failed: %w", err)
	}

	// Test Space-Track connection
	fmt.Println("Validating Space-Track credentials...")
	client, err := testSpaceTrackConnection(username, password)
	if err != nil {
		return fmt.Errorf("Space-Track connection test failed: %w", err)
	}
	_ = client // Client is validated, can be used later if needed

	// Test N2YO API connection
	fmt.Println("Validating N2YO API key...")
	if err := testN2YOConnection(apiKey); err != nil {
		return fmt.Errorf("N2YO API connection test failed: %w", err)
	}

	fmt.Println("All credentials validated successfully!")
	return nil
}

// testSpaceTrackConnection tests the Space-Track API connection.
func testSpaceTrackConnection(username, password string) (*http.Client, error) {
	// Import osint package functions - we'll need to make Login accessible or create a test function
	// For now, we'll create a minimal test here
	vals := url.Values{}
	vals.Add("identity", username)
	vals.Add("password", password)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.PostForm("https://www.space-track.org/ajaxauth/login", vals)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed (status: %d)", resp.StatusCode)
	}

	return client, nil
}

// testN2YOConnection tests the N2YO API connection with a simple request.
func testN2YOConnection(apiKey string) error {
	// Test with a simple request - get positions for ISS (NORAD ID 25544)
	// Using minimal parameters to reduce API usage
	testURL := fmt.Sprintf("https://api.n2yo.com/rest/v1/satellite/positions/25544/0/0/0/1/&apiKey=%s", apiKey)

	resp, err := http.Get(testURL)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try to read error message
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (status: %d): %s", resp.StatusCode, string(body))
	}

	// Verify response is valid JSON
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("invalid API response format: %w", err)
	}

	// Check for error in response
	if info, ok := result["info"].(map[string]interface{}); ok {
		if errMsg, ok := info["error"].(string); ok && errMsg != "" {
			return fmt.Errorf("API error: %s", errMsg)
		}
	}

	return nil
}

func main() {
	err := loadEnvFile()
	if err != nil {
		if err.Error() == ".env file not found" {
			fmt.Println("Note: .env file not found. Please provide credentials:")
		} else {
			fmt.Printf("Warning: Error loading .env file: %v\n", err)
			fmt.Println("Please provide credentials manually:")
		}
		fmt.Println()
	} else {
		fmt.Println("Loaded credentials from .env file")
	}

	checkEnvironmentalVariable("SPACE_TRACK_USERNAME")
	checkEnvironmentalVariable("SPACE_TRACK_PASSWORD")
	checkEnvironmentalVariable("N2YO_API_KEY")

	// Validate credentials format and test connections
	fmt.Println("\nValidating API credentials...")
	if err := validateCredentials(); err != nil {
		fmt.Printf("Warning: Credential validation failed: %v\n", err)
		fmt.Println("You may experience issues when using API features.")
		fmt.Println("Press Enter to continue anyway, or Ctrl+C to exit and fix credentials...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	cli.SatIntel()
}
