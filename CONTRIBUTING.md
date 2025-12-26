# Contributing to birdeye-go

Thank you for your interest in contributing to birdeye-go! This document provides guidelines for contributing to the project.

## Code of Conduct

Be respectful and constructive. We're all here to build good software.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/birdeye-go.git`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Run linter: `golangci-lint run`
7. Commit your changes with a clear message
8. Push to your fork
9. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Install dependencies

```bash
go mod download
```

### Run tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

### Run linter

```bash
golangci-lint run
```

## Code Standards

### General Guidelines

- Follow standard Go idioms and conventions
- Use `gofmt` to format code (or `goimports`)
- Write clear, descriptive variable and function names
- Keep functions focused and small
- Add comments for exported types and functions

### Error Handling

- Always check errors
- Wrap errors with context: `fmt.Errorf("failed to fetch price: %w", err)`
- Use the custom `APIError` type for API-related errors

### Testing

- Write tests for all new functionality
- Maintain 80%+ code coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

### Financial Precision

- **Always use `decimal.Decimal`** for financial values
- **Never use `float64`** for prices, amounts, or percentages
- This is non-negotiable for a financial API client

## Pull Request Guidelines

### Before Submitting

- [ ] Tests pass: `go test ./...`
- [ ] Linter passes: `golangci-lint run`
- [ ] Code is formatted: `gofmt -w .`
- [ ] Coverage is 80%+: `go test -cover ./...`
- [ ] Documentation is updated if needed

### PR Description

- Clearly describe what the PR does
- Reference any related issues
- Include examples if adding new features

### Commit Messages

Write clear, concise commit messages:

```
Add GetOHLCV method for candlestick data

- Add OHLCVData struct with open/high/low/close/volume
- Add GetOHLCV method with timeframe parameter
- Add unit tests with 95% coverage
```

## Adding New API Endpoints

When adding support for new Birdeye API endpoints:

1. **Create the response struct** in the appropriate file (or new file)
   - Use `decimal.Decimal` for all numeric financial values
   - Add JSON tags matching the API response
   - Add godoc comments explaining each field

2. **Add the client method**
   - Accept `context.Context` as first parameter
   - Validate required parameters
   - Use `parseResponse[T]` for type-safe parsing
   - Log with structured fields

3. **Write comprehensive tests**
   - Test success case
   - Test validation errors
   - Test 404 handling
   - Test success=false response
   - Test any special parsing (decimals, nested structs)

4. **Update documentation**
   - Add examples in godoc
   - Update README if it's a major feature

## Questions?

Open an issue for questions or discussions about larger changes before investing significant time.
