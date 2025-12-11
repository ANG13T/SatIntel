package osint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
)

// ErrorCode represents a unique error code for troubleshooting.
type ErrorCode string

const (
	// Authentication errors (1000-1099)
	ErrCodeAuthFailed          ErrorCode = "AUTH-1001"
	ErrCodeAuthCredentials     ErrorCode = "AUTH-1002"
	ErrCodeAuthConnection      ErrorCode = "AUTH-1003"
	ErrCodeAuthCookieJar       ErrorCode = "AUTH-1004"

	// API errors (1100-1199)
	ErrCodeAPIRequestFailed    ErrorCode = "API-1101"
	ErrCodeAPIResponseFailed   ErrorCode = "API-1102"
	ErrCodeAPIParseFailed      ErrorCode = "API-1103"
	ErrCodeAPINoData           ErrorCode = "API-1104"
	ErrCodeAPIInvalidEndpoint  ErrorCode = "API-1105"

	// Input validation errors (1200-1299)
	ErrCodeInputEmpty          ErrorCode = "INPUT-1201"
	ErrCodeInputInvalid        ErrorCode = "INPUT-1202"
	ErrCodeInputOutOfRange     ErrorCode = "INPUT-1203"
	ErrCodeInputFormat         ErrorCode = "INPUT-1204"

	// TLE errors (1300-1399)
	ErrCodeTLEInvalidFormat    ErrorCode = "TLE-1301"
	ErrCodeTLEParseFailed      ErrorCode = "TLE-1302"
	ErrCodeTLEInsufficientData ErrorCode = "TLE-1303"
	ErrCodeTLEChecksumFailed   ErrorCode = "TLE-1304"

	// File errors (1400-1499)
	ErrCodeFileNotFound        ErrorCode = "FILE-1401"
	ErrCodeFileReadFailed      ErrorCode = "FILE-1402"
	ErrCodeFilePathInvalid     ErrorCode = "FILE-1403"
	ErrCodeFilePermission      ErrorCode = "FILE-1404"

	// Satellite selection errors (1500-1599)
	ErrCodeSatNotFound         ErrorCode = "SAT-1501"
	ErrCodeSatInvalidNORAD     ErrorCode = "SAT-1502"
	ErrCodeSatNoResults        ErrorCode = "SAT-1503"

	// Network errors (1600-1699)
	ErrCodeNetworkTimeout      ErrorCode = "NET-1601"
	ErrCodeNetworkUnreachable  ErrorCode = "NET-1602"
	ErrCodeNetworkDNS          ErrorCode = "NET-1603"
)

// AppError represents a structured application error with code, message, and suggestions.
type AppError struct {
	Code       ErrorCode
	Message    string
	Context    string
	Suggestions []string
	OriginalErr error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.OriginalErr.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Display formats and displays the error with suggestions.
func (e *AppError) Display() {
	fmt.Println(color.Ize(color.Red, fmt.Sprintf("  [!] ERROR [%s]: %s", e.Code, e.Message)))
	
	if e.Context != "" {
		fmt.Println(color.Ize(color.Yellow, fmt.Sprintf("       Context: %s", e.Context)))
	}
	
	if len(e.Suggestions) > 0 {
		fmt.Println(color.Ize(color.Cyan, "       Suggestions:"))
		for i, suggestion := range e.Suggestions {
			fmt.Println(color.Ize(color.Cyan, fmt.Sprintf("         %d. %s", i+1, suggestion)))
		}
	}
	
	if e.OriginalErr != nil {
		fmt.Println(color.Ize(color.Gray, fmt.Sprintf("       Technical details: %v", e.OriginalErr)))
	}
}

// NewAppError creates a new AppError with the given code and message.
func NewAppError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:        code,
		Message:     message,
		Suggestions: getDefaultSuggestions(code),
	}
}

// NewAppErrorWithContext creates a new AppError with context.
func NewAppErrorWithContext(code ErrorCode, message, context string) *AppError {
	return &AppError{
		Code:        code,
		Message:     message,
		Context:     context,
		Suggestions: getDefaultSuggestions(code),
	}
}

// NewAppErrorWithErr wraps an original error in an AppError.
func NewAppErrorWithErr(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:        code,
		Message:     message,
		OriginalErr: err,
		Suggestions: getDefaultSuggestions(code),
	}
}

// getDefaultSuggestions returns default suggestions based on error code.
func getDefaultSuggestions(code ErrorCode) []string {
	suggestions := map[ErrorCode][]string{
		// Authentication errors
		ErrCodeAuthFailed: {
			"Verify your Space-Track credentials in the .env file",
			"Check if your account is active at space-track.org",
			"Ensure SPACE_TRACK_USERNAME and SPACE_TRACK_PASSWORD are set correctly",
		},
		ErrCodeAuthCredentials: {
			"Check your .env file for SPACE_TRACK_USERNAME and SPACE_TRACK_PASSWORD",
			"Verify credentials are not expired or changed",
			"Try logging in manually at space-track.org to verify your account",
		},
		ErrCodeAuthConnection: {
			"Check your internet connection",
			"Verify space-track.org is accessible",
			"Check if a firewall or proxy is blocking the connection",
			"Try again in a few moments - the service may be temporarily unavailable",
		},
		ErrCodeAuthCookieJar: {
			"Restart the application",
			"Check system permissions for cookie storage",
		},

		// API errors
		ErrCodeAPIRequestFailed: {
			"Check your internet connection",
			"Verify the API endpoint is correct",
			"Try again in a few moments",
			"Check if you're authenticated (run login first)",
		},
		ErrCodeAPIResponseFailed: {
			"The API returned an error response",
			"Check if your query parameters are valid",
			"Verify you have permission to access this data",
			"Try with different search criteria",
		},
		ErrCodeAPIParseFailed: {
			"The API response format may have changed",
			"Check if the data structure is correct",
			"Try refreshing the data",
		},
		ErrCodeAPINoData: {
			"Try adjusting your search criteria",
			"Verify the satellite NORAD ID exists",
			"Check if the satellite is still active",
		},
		ErrCodeAPIInvalidEndpoint: {
			"Verify the endpoint URL is correct",
			"Check API documentation for valid endpoints",
		},

		// Input validation errors
		ErrCodeInputEmpty: {
			"Please provide a value for this field",
			"Check if you accidentally pressed Enter without entering data",
		},
		ErrCodeInputInvalid: {
			"Enter a valid numeric value",
			"Remove any non-numeric characters (except decimal point and minus sign)",
			"Check the expected format for this input",
		},
		ErrCodeInputOutOfRange: {
			"Check the valid range for this input",
			"Latitude: -90 to 90",
			"Longitude: -180 to 180",
			"Altitude: 0 to 8848 meters (Mount Everest)",
		},
		ErrCodeInputFormat: {
			"Verify the input format matches the expected pattern",
			"Check for typos or extra characters",
		},

		// TLE errors
		ErrCodeTLEInvalidFormat: {
			"Ensure TLE data has exactly 2 lines",
			"Line 1 must start with '1 '",
			"Line 2 must start with '2 '",
			"Each line should be 69 characters long",
		},
		ErrCodeTLEParseFailed: {
			"Verify the TLE data is complete and not corrupted",
			"Check if all required fields are present",
			"Ensure the TLE format follows the standard specification",
		},
		ErrCodeTLEInsufficientData: {
			"Ensure you have both TLE lines",
			"Check if the data was truncated",
			"Verify the source of your TLE data",
		},
		ErrCodeTLEChecksumFailed: {
			"The TLE checksum validation failed",
			"Verify the TLE data is not corrupted",
			"Try fetching fresh TLE data from Space-Track",
		},

		// File errors
		ErrCodeFileNotFound: {
			"Verify the file path is correct",
			"Check if the file exists at the specified location",
			"Use an absolute path or ensure the file is in the current directory",
		},
		ErrCodeFileReadFailed: {
			"Check file permissions",
			"Ensure the file is not locked by another process",
			"Verify the file is not corrupted",
		},
		ErrCodeFilePathInvalid: {
			"Use a valid file path",
			"Avoid directory traversal patterns (../)",
			"Check for invalid characters in the path",
		},
		ErrCodeFilePermission: {
			"Check file and directory permissions",
			"Ensure you have read access to the file",
			"On Unix systems, check with: ls -l <file>",
		},

		// Satellite selection errors
		ErrCodeSatNotFound: {
			"Verify the NORAD ID is correct",
			"Try searching by satellite name instead",
			"Check if the satellite is still active",
		},
		ErrCodeSatInvalidNORAD: {
			"NORAD ID must be a numeric value",
			"Check for typos in the NORAD ID",
			"Verify the format: numeric only, no spaces",
		},
		ErrCodeSatNoResults: {
			"Try adjusting your search criteria",
			"Use broader search terms",
			"Check spelling of satellite name or country",
			"Try removing filters to see more results",
		},

		// Network errors
		ErrCodeNetworkTimeout: {
			"Check your internet connection",
			"The request took too long - try again",
			"Check if the API service is experiencing high load",
		},
		ErrCodeNetworkUnreachable: {
			"Verify your internet connection is active",
			"Check firewall settings",
			"Try accessing the API URL in a browser",
		},
		ErrCodeNetworkDNS: {
			"Check your DNS settings",
			"Verify you can resolve the API domain",
			"Try using a different DNS server (e.g., 8.8.8.8)",
		},
	}

	if sug, ok := suggestions[code]; ok {
		return sug
	}
	return []string{"Please check the error details and try again"}
}

// HandleError displays an error if it's an AppError, otherwise creates a generic error.
func HandleError(err error, defaultCode ErrorCode, defaultMessage string) {
	if err == nil {
		return
	}

	if appErr, ok := err.(*AppError); ok {
		appErr.Display()
		return
	}

	// Wrap generic errors
	appErr := NewAppErrorWithErr(defaultCode, defaultMessage, err)
	appErr.Display()
}

// HandleErrorWithContext displays an error with additional context.
func HandleErrorWithContext(err error, code ErrorCode, message, context string) {
	if err == nil {
		return
	}

	if appErr, ok := err.(*AppError); ok {
		appErr.Context = context
		appErr.Display()
		return
	}

	appErr := NewAppErrorWithContext(code, message, context)
	appErr.OriginalErr = err
	appErr.Display()
}

// ValidateInput checks if input is empty and returns an appropriate error.
func ValidateInput(value, fieldName string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return NewAppErrorWithContext(
			ErrCodeInputEmpty,
			fmt.Sprintf("%s cannot be empty", fieldName),
			fmt.Sprintf("Field: %s", fieldName),
		)
	}
	return nil
}

// cleanNumericInputForValidation removes non-numeric characters from input string.
func cleanNumericInputForValidation(input string) string {
	var result strings.Builder
	for _, char := range input {
		if (char >= '0' && char <= '9') || char == '.' || char == '-' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// ValidateNumericInput validates numeric input and returns an appropriate error.
func ValidateNumericInput(value, fieldName string, min, max float64) error {
	if err := ValidateInput(value, fieldName); err != nil {
		return err
	}

	// Clean the input
	cleaned := cleanNumericInputForValidation(value)
	if cleaned == "" {
		return NewAppErrorWithContext(
			ErrCodeInputInvalid,
			fmt.Sprintf("%s must be a valid number", fieldName),
			fmt.Sprintf("Field: %s, Value: %s", fieldName, value),
		)
	}

	// Parse and validate range if specified
	if min != 0 || max != 0 {
		var num float64
		var err error
		if strings.Contains(cleaned, ".") {
			num, err = strconv.ParseFloat(cleaned, 64)
		} else {
			var n int
			n, err = strconv.Atoi(cleaned)
			num = float64(n)
		}

		if err != nil {
			return NewAppErrorWithContext(
				ErrCodeInputInvalid,
				fmt.Sprintf("%s must be a valid number", fieldName),
				fmt.Sprintf("Field: %s, Value: %s", fieldName, value),
			)
		}

		if (min != 0 && num < min) || (max != 0 && num > max) {
			return NewAppErrorWithContext(
				ErrCodeInputOutOfRange,
				fmt.Sprintf("%s must be between %.2f and %.2f", fieldName, min, max),
				fmt.Sprintf("Field: %s, Value: %.2f", fieldName, num),
			)
		}
	}

	return nil
}

