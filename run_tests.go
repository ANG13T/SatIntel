// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TestModule represents a testable module in the application
type TestModule struct {
	Name        string
	Description string
	PackagePath string
	TestFile    string
}

// Available test modules
var testModules = []TestModule{
	{
		Name:        "main",
		Description: "Main package tests (env loading, password masking, credential management)",
		PackagePath: ".",
		TestFile:    "main_test.go",
	},
	{
		Name:        "osint",
		Description: "OSINT package tests (TLE parsing, API interactions, satellite data, export functionality)",
		PackagePath: "./osint",
		TestFile:    "osint_test.go",
	},
	{
		Name:        "cli",
		Description: "CLI package tests (menu handling, user input, display functions)",
		PackagePath: "./cli",
		TestFile:    "cli_test.go",
	},
	{
		Name:        "export",
		Description: "Export functionality tests (CSV, JSON, Text export for TLE, predictions, positions)",
		PackagePath: "./osint",
		TestFile:    "export_test.go",
	},
	// Add new modules here as you create them
}

// runTests executes go test for a specific module
func runTests(module TestModule, verbose bool, coverage bool, benchmark bool, testPattern string) error {
	args := []string{"test"}

	if verbose {
		args = append(args, "-v")
	}

	if coverage {
		args = append(args, "-cover")
	}

	if benchmark {
		args = append(args, "-bench=.", "-benchmem")
	}

	// Add test name pattern filter if specified
	if testPattern != "" {
		args = append(args, "-run", testPattern)
	}

	args = append(args, module.PackagePath)

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("\n%s Testing module: %s %s\n", strings.Repeat("=", 40), module.Name, strings.Repeat("=", 40))
	fmt.Printf("Description: %s\n", module.Description)
	fmt.Printf("Package: %s\n", module.PackagePath)
	if testPattern != "" {
		fmt.Printf("Test Pattern: %s\n", testPattern)
	}
	fmt.Printf("Command: go %s\n\n", strings.Join(args, " "))

	return cmd.Run()
}

// listModules displays all available test modules
func listModules() {
	fmt.Println("\nAvailable Test Modules:")
	fmt.Println(strings.Repeat("-", 60))
	for i, module := range testModules {
		fmt.Printf("%d. %s\n", i+1, module.Name)
		fmt.Printf("   Description: %s\n", module.Description)
		fmt.Printf("   Package: %s\n", module.PackagePath)
		fmt.Printf("   Test File: %s\n\n", module.TestFile)
	}
}

// checkTestFiles verifies which test files exist
func checkTestFiles() {
	fmt.Println("\nTest File Status:")
	fmt.Println(strings.Repeat("-", 60))
	for _, module := range testModules {
		testPath := filepath.Join(module.PackagePath, module.TestFile)
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			fmt.Printf("❌ %s - Missing: %s\n", module.Name, testPath)
		} else {
			fmt.Printf("✅ %s - Exists: %s\n", module.Name, testPath)
		}
	}
}

func main() {
	var (
		moduleFlag    = flag.String("module", "", "Run tests for specific module (use 'list' to see available modules)")
		allFlag       = flag.Bool("all", false, "Run all test modules")
		verboseFlag   = flag.Bool("v", false, "Verbose output")
		coverageFlag  = flag.Bool("cover", false, "Show coverage information")
		benchFlag     = flag.Bool("bench", false, "Run benchmarks")
		checkFlag     = flag.Bool("check", false, "Check which test files exist")
		helpFlag      = flag.Bool("help", false, "Show help message")
		runFlag       = flag.String("run", "", "Run only tests matching the pattern (e.g., 'TestExport' for export tests)")
	)

	flag.Parse()

	if *helpFlag {
		fmt.Println("SatIntel Test Runner")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("\nUsage:")
		fmt.Println("  go run run_tests.go [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  go run run_tests.go -all                    # Run all tests")
		fmt.Println("  go run run_tests.go -module=main            # Run main package tests")
		fmt.Println("  go run run_tests.go -all -cover             # Run all tests with coverage")
		fmt.Println("  go run run_tests.go -module=osint -v        # Run osint tests verbosely")
		fmt.Println("  go run run_tests.go -module=export -run TestExport  # Run only export tests")
		fmt.Println("  go run run_tests.go -check                  # Check test file status")
		fmt.Println("  go run run_tests.go -module=list            # List available modules")
		return
	}

	if *checkFlag {
		checkTestFiles()
		return
	}

	if *moduleFlag == "list" {
		listModules()
		return
	}

	// Find specific module if requested
	if *moduleFlag != "" {
		var foundModule *TestModule
		for i := range testModules {
			if testModules[i].Name == *moduleFlag {
				foundModule = &testModules[i]
				break
			}
		}

		if foundModule == nil {
			fmt.Printf("Error: Module '%s' not found.\n", *moduleFlag)
			fmt.Println("\nAvailable modules:")
			for _, m := range testModules {
				fmt.Printf("  - %s\n", m.Name)
			}
			os.Exit(1)
		}

		if err := runTests(*foundModule, *verboseFlag, *coverageFlag, *benchFlag, *runFlag); err != nil {
			os.Exit(1)
		}
		return
	}

	// Run all modules if -all flag is set, otherwise show help
	if *allFlag {
		fmt.Println("SatIntel Test Runner - Running All Modules")
		fmt.Println(strings.Repeat("=", 60))

		failed := false
		for _, module := range testModules {
			// Check if test file exists before running
			testPath := filepath.Join(module.PackagePath, module.TestFile)
			if _, err := os.Stat(testPath); os.IsNotExist(err) {
				fmt.Printf("\n⚠️  Skipping %s - test file not found: %s\n", module.Name, testPath)
				continue
			}

			if err := runTests(module, *verboseFlag, *coverageFlag, *benchFlag, *runFlag); err != nil {
				failed = true
				fmt.Printf("\n❌ Tests failed for module: %s\n", module.Name)
			} else {
				fmt.Printf("\n✅ Tests passed for module: %s\n", module.Name)
			}
		}

		if failed {
			fmt.Println("\n" + strings.Repeat("=", 60))
			fmt.Println("Some tests failed. See output above for details.")
			os.Exit(1)
		} else {
			fmt.Println("\n" + strings.Repeat("=", 60))
			fmt.Println("All tests passed! ✅")
		}
		return
	}

	// Default: show help
	fmt.Println("SatIntel Test Runner")
	fmt.Println("Use -help to see usage information")
	fmt.Println("Use -module=list to see available modules")
	fmt.Println("Use -all to run all tests")
	fmt.Println("Use -check to check test file status")
}

