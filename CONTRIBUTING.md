# Contributing to Evertec/paySmart Processadora SDK for Go

Thank you for your interest in contributing! This document provides guidelines for contributing to the SDK.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/evertec-golang-sdk-processadora.git
   cd evertec-golang-sdk-processadora
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/henriqueatila/evertec-golang-sdk-processadora.git
   ```

## Development Setup

### Requirements

- Go 1.21 or later
- golangci-lint (for linting)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Running Linter

```bash
golangci-lint run
```

## Making Changes

### Branch Naming

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring

### Commit Messages

Follow conventional commits format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or modifications
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

Examples:
```
feat(client): add support for batch card operations
fix(authorization): handle nil pointer in purchase handler
docs(readme): add webhook configuration example
test(accounts): add integration tests for limits
```

## Code Style

### Go Guidelines

1. **Follow Go conventions**: Use `gofmt` and `go vet`
2. **Error handling**: Always handle errors explicitly
3. **Documentation**: Add GoDoc comments for exported functions and types
4. **Naming**: Use clear, descriptive names following Go conventions

### Example Code Style

```go
// GetAccount retrieves account details by ID.
// Returns an error if the account is not found.
func (c *Client) GetAccount(ctx context.Context, accountID string) (*types.Account, error) {
    if accountID == "" {
        return nil, errors.New("accountID is required")
    }

    path := fmt.Sprintf("/accounts/%s", accountID)

    var resp types.Account
    if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
        return nil, fmt.Errorf("get account: %w", err)
    }

    return &resp, nil
}
```

### Testing Guidelines

1. **Table-driven tests**: Use table-driven tests for multiple scenarios
2. **Parallel tests**: Use `t.Parallel()` when tests are independent
3. **Mock servers**: Use `httptest.NewServer` for HTTP mocking
4. **Coverage**: Maintain 80%+ test coverage

Example test:

```go
func TestClient_GetAccount(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        accountID  string
        mockStatus int
        mockResp   interface{}
        wantErr    bool
    }{
        {
            name:       "success",
            accountID:  "acc-001",
            mockStatus: http.StatusOK,
            mockResp:   &types.Account{AccountID: "acc-001"},
            wantErr:    false,
        },
        {
            name:       "not found",
            accountID:  "acc-999",
            mockStatus: http.StatusNotFound,
            mockResp:   map[string]string{"error": "not found"},
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // test implementation
        })
    }
}
```

## Pull Request Process

1. **Create feature branch** from `main`
2. **Make changes** following the guidelines above
3. **Write/update tests** for your changes
4. **Run tests and linter** to ensure everything passes
5. **Update documentation** if needed
6. **Submit PR** with clear description

### PR Checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Coverage maintained or improved
- [ ] Documentation updated
- [ ] Commit messages follow convention

## Reporting Issues

When reporting issues, please include:

1. **Go version**: `go version`
2. **SDK version**: Check your `go.mod`
3. **Steps to reproduce**: Clear, numbered steps
4. **Expected behavior**: What should happen
5. **Actual behavior**: What actually happens
6. **Error messages**: Full error output if applicable

## Security

If you discover a security vulnerability, please do NOT open a public issue. Instead, email the maintainers directly.

## Questions?

Feel free to open an issue for questions about contributing.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
