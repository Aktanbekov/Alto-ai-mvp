# Test Summary

This document describes the automated tests for the AltoAI MVP project.

## Test Coverage

### Interview Package (`interview/`)
- **analyzer_test.go**: Tests for scoring system, grade calculation, percentage conversion, and recommendations
- **session_test.go**: Tests for session creation, storage, and retrieval
- **models_test.go**: Tests for JSON serialization/deserialization of analysis models
- **evaluation_test.go**: Tests for extracting strengths, weaknesses, red flags, and session summaries
- **questions_test.go**: Tests for loading and validating interview questions

### Internal Services (`internal/services/`)
- **auth_service_test.go**: Tests for password hashing, comparison, and email code generation

### Internal Middleware (`internal/middleware/`)
- **cors_test.go**: Tests for CORS middleware functionality

### Package Response (`pkg/response/`)
- **response_test.go**: Tests for HTTP response helpers

## Running Tests

### Run all tests:
```bash
go test ./interview/... ./internal/... ./pkg/... -v
```

### Run specific package tests:
```bash
go test ./interview/... -v
go test ./internal/services/... -v
go test ./internal/middleware/... -v
```

### Run tests before starting servers:
```bash
./run.sh
```

The `run.sh` script automatically runs all tests before starting the backend and frontend servers. If any test fails, the servers will not start.

## Test Statistics

- **Total Test Files**: 7
- **Test Functions**: ~25+
- **Coverage Areas**:
  - Scoring and grading system
  - Session management
  - Authentication services
  - CORS middleware
  - JSON serialization
  - Helper functions

## Notes

- The `ai_tests` package is excluded from automated testing as it contains incomplete test code
- Tests use Go's standard `testing` package
- All tests should pass before deploying or starting servers


