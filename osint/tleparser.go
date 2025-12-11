package osint

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"
)

// TLEParser provides an interactive menu for parsing TLE data from different sources.
func TLEParser() {
	options, _ := os.ReadFile("txt/tle_parser.txt")
	opt, _ := gradient.NewGradient("#1179ef", "cyan")
	opt.Print("\n" + string(options))
	var selection int = Option(0, 3)

	if selection == 1 {
		TLETextFile()
	} else if selection == 2 {
		TLEPlainString()
	}
}

// validateFilePath checks if a file path is safe and valid.
// It prevents directory traversal attacks and validates path format.
func validateFilePath(path string) error {
	// Trim whitespace
	path = strings.TrimSpace(path)

	// Check for empty path
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check for null bytes (potential injection)
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("file path contains invalid characters")
	}

	// Check for directory traversal patterns
	dangerousPatterns := []string{
		"..",
		"../",
		"..\\",
		"/..",
		"\\..",
		"//",
		"\\\\",
	}
	pathNormalized := strings.ToLower(path)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(pathNormalized, strings.ToLower(pattern)) {
			return fmt.Errorf("directory traversal detected in path")
		}
	}

	// Check path length (reasonable limit)
	if len(path) > 4096 {
		return fmt.Errorf("file path is too long")
	}

	// Clean the path to resolve any remaining issues
	cleanedPath := filepath.Clean(path)

	// Check if cleaned path still contains dangerous patterns
	if strings.Contains(cleanedPath, "..") {
		return fmt.Errorf("invalid file path")
	}

	// Check if path is absolute and potentially dangerous
	// (Optional: you might want to restrict to relative paths only)
	// For now, we'll allow absolute paths but validate them

	return nil
}

// TLETextFile reads TLE data from a text file and parses it.
func TLETextFile() {
	fmt.Print("\n ENTER TEXT FILE PATH > ")
	var path string
	fmt.Scanln(&path)

	// Validate file path before attempting to open
	if err := validateFilePath(path); err != nil {
		appErr := NewAppErrorWithContext(ErrCodeFilePathInvalid, "Invalid file path", fmt.Sprintf("Path: %s", path))
		appErr.OriginalErr = err
		appErr.Display()
		return
	}

	// Clean the path after validation
	path = filepath.Clean(strings.TrimSpace(path))

	// Check if file exists and is a regular file (not a directory)
	fileInfo, err := os.Stat(path)
	if err != nil {
		context := fmt.Sprintf("File path: %s", path)
		HandleErrorWithContext(err, ErrCodeFileNotFound, "Failed to access TLE file", context)
		return
	}

	// Ensure it's a file, not a directory
	if fileInfo.IsDir() {
		err := NewAppErrorWithContext(ErrCodeFilePathInvalid, "Path is a directory, not a file", fmt.Sprintf("Path: %s", path))
		err.Display()
		return
	}

	file, err := os.Open(path)
	if err != nil {
		context := fmt.Sprintf("File path: %s", path)
		HandleErrorWithContext(err, ErrCodeFileReadFailed, "Failed to open TLE file", context)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	var count int = 0

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
		count += 1
	}

	if count < 2 || count > 3 {
		err := NewAppErrorWithContext(
			ErrCodeTLEInvalidFormat,
			"Invalid TLE format - file must contain 2 or 3 lines",
			fmt.Sprintf("File: %s, Lines found: %d", path, count),
		)
		err.Display()
		return
	}

	var output TLE
	var lineOne, lineTwo string

	if count == 3 {
		var satelliteName string = txtlines[0]
		lineOne = txtlines[1]
		lineTwo = txtlines[2]
		output = ConstructTLE(satelliteName, lineOne, lineTwo)
	} else {
		lineOne = txtlines[0]
		lineTwo = txtlines[1]
		output = ConstructTLE("UNSPECIFIED", lineOne, lineTwo)
	}

	// Validate TLE parsing before displaying
	parsingFailed := false
	line1Fields := strings.Fields(lineOne)
	line2Fields := strings.Fields(lineTwo)

	if len(line1Fields) < 4 || len(line2Fields) < 3 {
		parsingFailed = true
	} else if output.SatelliteCatalogNumber == 0 && output.InternationalDesignator == "" && output.ElementSetEpoch == 0.0 {
		parsingFailed = true
	}

	if parsingFailed {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to parse TLE data"))
		fmt.Println(color.Ize(color.Red, fmt.Sprintf("       Line 1 fields: %d (minimum required: 4)", len(line1Fields))))
		fmt.Println(color.Ize(color.Red, fmt.Sprintf("       Line 2 fields: %d (minimum required: 3)", len(line2Fields))))
		if len(line1Fields) >= 4 && len(line2Fields) >= 3 {
			fmt.Println(color.Ize(color.Red, "       Note: Field count is sufficient, but parsing failed. Check TLE format."))
		}
		return
	}

	PrintTLE(output)
}

// TLEPlainString prompts the user to enter TLE data line by line and parses it.
func TLEPlainString() {
	scanner := bufio.NewScanner(os.Stdin)
	var lineOne string
	var lineTwo string
	var lineThree string
	fmt.Print("\n ENTER LINE ONE (leave blank for unspecified name)  >  ")
	scanner.Scan()
	lineOne = scanner.Text()

	fmt.Print("\n ENTER LINE TWO  >  ")
	scanner.Scan()
	lineTwo = scanner.Text()

	fmt.Print("\n ENTER LINE THREE  >  ")
	scanner.Scan()
	lineThree = scanner.Text()

	if lineOne == "" {
		lineOne = "UNSPECIFIED"
	}

	output := ConstructTLE(lineOne, lineTwo, lineThree)

	// Validate TLE parsing before displaying
	parsingFailed := false
	line1Fields := strings.Fields(lineTwo)
	line2Fields := strings.Fields(lineThree)

	if len(line1Fields) < 4 || len(line2Fields) < 3 {
		parsingFailed = true
	} else if output.SatelliteCatalogNumber == 0 && output.InternationalDesignator == "" && output.ElementSetEpoch == 0.0 {
		parsingFailed = true
	}

	if parsingFailed {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to parse TLE data"))
		fmt.Println(color.Ize(color.Red, fmt.Sprintf("       Line 1 fields: %d (minimum required: 4)", len(line1Fields))))
		fmt.Println(color.Ize(color.Red, fmt.Sprintf("       Line 2 fields: %d (minimum required: 3)", len(line2Fields))))
		if len(line1Fields) >= 4 && len(line2Fields) >= 3 {
			fmt.Println(color.Ize(color.Red, "       Note: Field count is sufficient, but parsing failed. Check TLE format."))
		}
		return
	}

	PrintTLE(output)
}
