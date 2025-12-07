# Contributing to Driftlock

Thank you for your interest in contributing to Driftlock! We welcome contributions from the community.

## Getting Started

### Prerequisites

- Go 1.22+ for backend development
- Node.js 18+ for frontend development
- Rust 1.70+ for CBAD core development
- Docker and Docker Compose for local development
- PostgreSQL 15+ (can run via Docker)

### Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/your-username/driftlock.git
   cd driftlock
   ```

2. **Set up Go modules**
   ```bash
   cd collector-processor
   go mod download
   ```

3. **Set up the frontend**
   ```bash
   cd landing-page
   npm install
   ```

4. **Set up CBAD core (Rust)**
   ```bash
   cd cbad-core
   cargo build
   ```

5. **Start local services**
   ```bash
   # From the root directory
   docker compose up -d postgres redis
   ```

6. **Run database migrations**
   ```bash
   cd collector-processor
   goose -dir api/migrations postgres "$DATABASE_URL" up
   ```

### Running Tests

```bash
# Go backend tests
go test ./... -v

# Frontend tests
cd landing-page
npm test

# Rust tests
cd cbad-core
cargo test
```

## Code Style

### Go

- Follow the standard Go formatting conventions
- Run `make fmt` to format code
- Run `make lint` to check for linting issues
- Use `golangci-lint` for comprehensive linting

### Vue/TypeScript

- Use the provided ESLint and Prettier configurations
- Run `npm run lint` and `npm run format`
- Follow Vue 3 Composition API patterns
- Use TypeScript for type safety

### Rust

- Use `rustfmt` for formatting
- Use `clippy` for linting
- Follow Rust idioms and conventions

## Pull Request Process

1. Create a new branch from `main`
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes with atomic, well-documented commits

3. Ensure all tests pass and code is properly formatted

4. Write a clear pull request description:
   - What problem does this solve?
   - How did you solve it?
   - Are there any breaking changes?
   - Include screenshots for UI changes

5. Link to any relevant Linear issues

6. Request review from the maintainers

## Development Commands

```bash
# Start the API server
make run
# or
cd collector-processor && go run ./cmd/driftlock-http

# Start the frontend
cd landing-page && npm run dev

# Build all components
make build

# Run all tests
make test

# Format code
make fmt

# Run linters
make lint
```

## Project Structure

```
driftlock/
├── cbad-core/                   # Rust compression algorithms
├── collector-processor/        # Go API server and processor
│   ├── cmd/driftlock-http/     # HTTP API server
│   ├── driftlockcbad/         # CBAD integration
│   └── internal/             # Internal packages
├── landing-page/              # Vue.js frontend
├── api/                      # API definitions and migrations
├── docs/                     # Documentation
├── functions/                # Firebase Functions
├── extensions/               # VS Code extension
└── scripts/                  # Utility scripts
```

## Reporting Issues

- Use [GitHub Issues](https://github.com/Shannon-Labs/driftlock/issues) for bug reports
- Use [Linear](https://linear.app/shannon-labs) for feature requests (internal team)
- Include detailed reproduction steps for bugs
- Provide logs and error messages when possible

## Security

If you discover a security vulnerability, please report it privately to security@shannonlabs.dev. Do not open a public issue.

## License

By contributing to Driftlock, you agree that your contributions will be licensed under the same license as the project (Apache 2.0 or Commercial).

## Getting Help

- Check the [documentation](docs/README.md)
- Look at existing issues and pull requests
- Reach out to the maintainers
- Join our [Discord community](https://discord.gg/driftlock)

## Development Tips

- Use the provided VS Code extension for better development experience
- The demo endpoint (`/v1/demo/detect`) is great for testing
- Check `docs/AI_ROUTING.md` for guidance on which agents to use for different tasks
- Refer to `CLAUDE.md` for quick reference commands

Thank you for contributing! 🚀