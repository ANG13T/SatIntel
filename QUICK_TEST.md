# Quick Test Reference

## Run All Tests
```bash
go run run_tests.go -all
```

## Run Specific Module
```bash
go run run_tests.go -module=main
go run run_tests.go -module=osint
go run run_tests.go -module=cli
```

## Common Commands

| Command | What It Does |
|---------|--------------|
| `go run run_tests.go -all` | Run all test modules |
| `go run run_tests.go -all -cover` | Run all tests with coverage |
| `go run run_tests.go -module=main -v` | Run main tests verbosely |
| `go run run_tests.go -module=list` | List available modules |
| `go run run_tests.go -check` | Check test file status |
| `go run run_tests.go -help` | Show full help |

## Adding New Test Module

1. Create test file: `newpackage/newpackage_test.go`
2. Add to `run_tests.go`:
```go
{
    Name:        "newpackage",
    Description: "Your description",
    PackagePath: "./newpackage",
    TestFile:    "newpackage_test.go",
},
```
3. Run: `go run run_tests.go -module=newpackage`

## Direct Go Test (Alternative)

```bash
go test ./...              # All tests
go test -cover ./...       # With coverage
go test -v ./main          # Specific package
```





