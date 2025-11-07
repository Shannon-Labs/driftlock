# Contributing to Shannon Labs DriftLock

Thank you for your interest in contributing to DriftLock! This document provides guidelines and information about contributing to this project.

## Code of Conduct

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md). Please read and follow it in all your interactions with the project.

## Getting Started

### Prerequisites

- Go 1.24+ (for API server)
- Rust 1.70+ (for anomaly detection core)
- Node.js 18+ (for dashboard)
- Docker and Docker Compose
- Git

### Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/Shannon-Labs/driftlock.git
   cd driftlock
   ```

2. **Install dependencies**
   ```bash
   # Rust dependencies
   cd src/anomaly-detection
   cargo build

   # Go dependencies
   cd ../../src/api-server
   go mod download

   # Node.js dependencies
   cd ../../src/dashboard
   npm install
   ```

3. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the development environment**
   ```bash
   docker-compose up -d
   ```

## How to Contribute

### Reporting Issues

- Search existing issues before creating a new one
- Use the provided issue templates
- Provide clear, reproducible steps for bugs
- Include environment details (OS, versions, etc.)

### Submitting Pull Requests

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow the existing code style
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   # Run Rust tests
   cd src/anomaly-detection && cargo test

   # Run Go tests
   cd ../api-server && go test ./...

   # Run Node.js tests
   cd ../dashboard && npm test
   ```

4. **Commit your changes**
   ```bash
   git commit -m "feat: add your feature description"
   ```

5. **Push and create a pull request**
   ```bash
   git push origin feature/your-feature-name
   ```

### Code Style Guidelines

#### Rust
- Use `rustfmt` for formatting
- Use `clippy` for linting
- Follow Rust naming conventions

#### Go
- Use `gofmt` for formatting
- Follow Go conventions and idioms
- Write clear, commented code

#### TypeScript/React
- Use ESLint and Prettier
- Follow the existing component structure
- Use TypeScript for all new code

### Documentation

- Update README.md if adding new features
- Add code comments for complex logic
- Update API documentation for API changes
- Include examples in documentation

## Development Workflow

### Branch Structure

- `main`: Stable, production-ready code
- `develop`: Integration branch for features
- `feature/*`: Feature branches
- `bugfix/*`: Bug fix branches
- `release/*`: Release preparation branches

### Release Process

1. Update version numbers
2. Update CHANGELOG.md
3. Create release tag
4. Build and publish release artifacts

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific component tests
cd src/anomaly-detection && cargo test
cd src/api-server && go test ./...
cd src/dashboard && npm test
```

### Test Coverage

- Aim for >80% code coverage for new code
- Add integration tests for complex workflows
- Include performance benchmarks for critical paths

## Getting Help

- Check [documentation](docs/)
- Search [existing issues](https://github.com/Shannon-Labs/driftlock/issues)
- Start a [discussion](https://github.com/Shannon-Labs/driftlock/discussions)
- Contact maintainers at [Shannon-Labs@example.com]

## Security

If you discover a security vulnerability, please follow our [security policy](SECURITY.md) and report it privately.

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.