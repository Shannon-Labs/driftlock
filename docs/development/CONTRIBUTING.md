# Contributing to Shannon Labs Driftlock

Thank you for your interest in contributing to Driftlock! This document provides guidelines and information about contributing to this project.

## Code of Conduct

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md). Please read and follow it in all your interactions with the project.

## Getting Started

### Prerequisites

- **Rust** 1.75+ (for the API server and CBAD core)
- **Node.js** 18+ (for the dashboard)
- **Docker** and Docker Compose
- **PostgreSQL** 15+ (or use Docker)
- **Git**

### Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/Shannon-Labs/driftlock.git
   cd driftlock
   ```

2. **Build the project**
   ```bash
   # Build API server
   cargo build -p driftlock-api --release

   # Build CBAD core
   cargo build -p cbad-core --release

   # Node.js dependencies (for dashboard)
   cd landing-page
   npm install
   ```

3. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start PostgreSQL**
   ```bash
   docker run --name driftlock-postgres \
     -e POSTGRES_DB=driftlock \
     -e POSTGRES_USER=driftlock \
     -e POSTGRES_PASSWORD=driftlock \
     -p 5432:5432 \
     -d postgres:15
   ```

5. **Run the development environment**
   ```bash
   DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
     cargo run -p driftlock-api
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
   # Run all tests
   cargo test --workspace

   # Run API tests
   cargo test -p driftlock-api

   # Run CBAD tests
   cargo test -p cbad-core --release

   # Run frontend build
   cd landing-page && npm run build
   ```

4. **Format and lint**
   ```bash
   # Format Rust code
   cargo fmt --all

   # Lint Rust code
   cargo clippy --all-targets -- -D warnings
   ```

5. **Commit your changes**
   ```bash
   git commit -m "feat: add your feature description"
   ```

6. **Push and create a pull request**
   ```bash
   git push origin feature/your-feature-name
   ```

### Code Style Guidelines

#### Rust
- Use `cargo fmt` for formatting
- Use `cargo clippy` for linting
- Follow Rust naming conventions
- Write clear, documented code
- Use `tracing` for logging (we call it "driftlog")

#### TypeScript/Vue
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
# Run all workspace tests
cargo test --workspace

# Run specific crate tests
cargo test -p driftlock-api
cargo test -p cbad-core --release
cargo test -p driftlock-db

# Run with output
cargo test --workspace -- --nocapture

# Run single test
cargo test test_health_check

# Run frontend build verification
cd landing-page && npm run build
```

### Test Coverage

- Aim for >80% code coverage for new code
- Add integration tests for complex workflows
- Include performance benchmarks for critical paths

## Project Structure

```
driftlock/
├── Cargo.toml              # Workspace root
├── cbad-core/              # CBAD detection algorithms
│   └── src/
├── crates/
│   ├── driftlock-api/      # Axum HTTP server
│   ├── driftlock-db/       # Database layer (sqlx)
│   ├── driftlock-auth/     # Firebase + API key auth
│   ├── driftlock-billing/  # Stripe integration
│   └── driftlock-email/    # SendGrid emails
├── landing-page/           # Vue 3 dashboard
└── docs/                   # Documentation
```

## Getting Help

- Check [documentation](docs/)
- Search [existing issues](https://github.com/Shannon-Labs/driftlock/issues)
- Start a [discussion](https://github.com/Shannon-Labs/driftlock/discussions)
- Contact maintainers at [hunter@shannonlabs.dev](mailto:hunter@shannonlabs.dev)

## Security

If you discover a security vulnerability, please follow our [security policy](SECURITY.md) and report it privately.

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
