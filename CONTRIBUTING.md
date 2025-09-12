# Contributing to Go Micro Framework

Thank you for your interest in contributing to Go Micro Framework! This document provides guidelines and information for contributors.

## 🎯 How to Contribute

### Types of Contributions

We welcome several types of contributions:

- **Bug Reports**: Help us identify and fix issues
- **Feature Requests**: Suggest new features or improvements
- **Code Contributions**: Submit code changes, bug fixes, or new features
- **Documentation**: Improve or add documentation
- **Examples**: Create example projects or tutorials
- **Testing**: Add tests or improve test coverage

### Getting Started

1. **Fork the Repository**
   ```bash
   git clone https://github.com/your-username/go-micro-framework.git
   cd go-micro-framework
   ```

2. **Set Up Development Environment**
   ```bash
   # Install dependencies
   go mod tidy
   
   # Install development tools
   make tools
   
   # Run tests to ensure everything works
   make test
   ```

3. **Create a Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b bugfix/issue-number
   ```

## 🛠️ Development Guidelines

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Use `golangci-lint` for linting
- Write comprehensive tests for new features
- Document public APIs with Go doc comments

### Commit Message Format

We use conventional commits format:

```
type(scope): description

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(cli): add new command for service validation
fix(generator): resolve template rendering issue
docs(readme): update installation instructions
test(core): add unit tests for bootstrap engine
```

### Testing Requirements

- All new features must include tests
- Aim for >80% code coverage
- Include unit tests, integration tests, and E2E tests where appropriate
- Tests should be deterministic and not depend on external services

### Documentation Requirements

- Update README.md if you add new features
- Add or update API documentation
- Include examples in your pull request
- Update CHANGELOG.md for significant changes

## 🚀 Pull Request Process

### Before Submitting

1. **Ensure Tests Pass**
   ```bash
   make test
   make lint
   make security
   ```

2. **Update Documentation**
   - Update README.md if needed
   - Add/update API documentation
   - Include examples

3. **Update CHANGELOG.md**
   - Add your changes to the appropriate section
   - Follow the existing format

### Submitting a Pull Request

1. **Create Pull Request**
   - Use a clear, descriptive title
   - Provide a detailed description
   - Link any related issues

2. **Pull Request Template**
   ```markdown
   ## Description
   Brief description of changes

   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update

   ## Testing
   - [ ] Unit tests pass
   - [ ] Integration tests pass
   - [ ] Manual testing completed

   ## Checklist
   - [ ] Code follows style guidelines
   - [ ] Self-review completed
   - [ ] Documentation updated
   - [ ] CHANGELOG.md updated
   ```

3. **Review Process**
   - Maintainers will review your PR
   - Address any feedback promptly
   - Be responsive to questions or requests for changes

## 🏗️ Project Structure

Understanding the project structure will help you contribute effectively:

```
go-micro-framework/
├── cmd/                    # CLI application
│   └── microframework/
├── internal/               # Internal packages
│   ├── core/              # Core framework logic
│   ├── generator/         # Code generation
│   ├── templates/         # Go templates
│   └── validators/        # Configuration validation
├── pkg/                   # Public packages
├── scripts/               # Build and deployment scripts
├── docs/                  # Documentation
├── examples/              # Example projects
└── tests/                 # Test files
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test types
make test-unit
make test-integration
make test-e2e

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test ./internal/core/...
```

### Writing Tests

- Use table-driven tests where appropriate
- Mock external dependencies
- Test both success and error cases
- Use descriptive test names

**Example:**
```go
func TestServiceGenerator_GenerateService(t *testing.T) {
    tests := []struct {
        name    string
        config  *ServiceConfig
        wantErr bool
    }{
        {
            name: "valid config",
            config: &ServiceConfig{
                Name: "test-service",
                Type: "rest",
            },
            wantErr: false,
        },
        {
            name:    "invalid config",
            config:  nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            generator := NewServiceGenerator(logger)
            err := generator.GenerateService(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateService() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 📚 Documentation

### Code Documentation

- Use Go doc comments for public APIs
- Include examples in documentation
- Document complex algorithms or business logic

**Example:**
```go
// ServiceGenerator generates microservices using templates and configuration.
// It supports various service types including REST, gRPC, and event-driven services.
//
// Example:
//   generator := NewServiceGenerator(logger)
//   config := &ServiceConfig{
//       Name: "user-service",
//       Type: "rest",
//   }
//   err := generator.GenerateService(config)
func NewServiceGenerator(logger *logrus.Logger) *ServiceGenerator {
    // implementation
}
```

### README Updates

When adding new features:
- Update the features list
- Add usage examples
- Update installation instructions if needed
- Add configuration examples

## 🐛 Bug Reports

### Before Reporting

1. Check existing issues
2. Ensure you're using the latest version
3. Try to reproduce the issue

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Run command '...'
2. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g. Ubuntu 20.04]
- Go version: [e.g. 1.21]
- Framework version: [e.g. 1.0.0]

**Additional context**
Add any other context about the problem here.
```

## 💡 Feature Requests

### Before Requesting

1. Check existing feature requests
2. Consider if it fits the project's scope
3. Think about implementation approach

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Alternative solutions or workarounds.

**Additional context**
Add any other context or screenshots about the feature request.
```

## 🔧 Development Tools

### Required Tools

- Go 1.21+
- Docker (for testing)
- Make

### Recommended Tools

- VS Code with Go extension
- Git hooks for formatting
- Pre-commit hooks

### VS Code Settings

```json
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "go.coverOnSave": true
}
```

## 📋 Code Review Checklist

### For Contributors

- [ ] Code follows style guidelines
- [ ] Tests are included and pass
- [ ] Documentation is updated
- [ ] Commit messages follow convention
- [ ] No sensitive information in code
- [ ] Error handling is appropriate
- [ ] Performance considerations addressed

### For Reviewers

- [ ] Code quality is good
- [ ] Tests are comprehensive
- [ ] Documentation is clear
- [ ] No security issues
- [ ] Performance is acceptable
- [ ] Breaking changes are documented

## 🚀 Release Process

### Version Numbering

We use semantic versioning (SemVer):
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

### Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated
- [ ] Version is bumped
- [ ] Release notes are written
- [ ] GitHub release is created

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Respect different viewpoints

### Getting Help

- Check documentation first
- Search existing issues
- Ask questions in discussions
- Join our Discord community

## 📞 Contact

- **GitHub Issues**: For bugs and feature requests
- **Discussions**: For questions and ideas
- **Discord**: For real-time chat
- **Email**: For security issues

## 🙏 Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation
- Community highlights

Thank you for contributing to Go Micro Framework! 🚀
