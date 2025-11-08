# Contributing Guide

We welcome contributions that reinforce the deterministic, explainable goals of Driftlock. To keep the repository audit-friendly, please follow the guidelines below.

## Commit Message Convention
Adopt a lightweight fork of Conventional Commits:
```
<type>(<scope>): <subject>
```
- **type**: `feat`, `fix`, `docs`, `test`, `chore`, `refactor`, `perf`, `build`, `ci`.
- **scope**: optional but encouraged (e.g., `cbad-core`, `router`, `api`, `ui`).
- **subject**: concise, imperative, ≤72 characters.

Include a body when explaining decisions, referencing decision log entries, or linking to compliance requirements.

### Examples
```
feat(cbad-core): add normalized compression distance metric

Implements NCD calculation for baseline/window comparison as outlined
in ALGORITHMS.md. Uses deterministic compression with configurable
backends (zstd, lz4, gzip).

Refs: decision-log.md entry 2025-01-09 compression abstraction
```

```
fix(api): handle malformed JSON in event ingestion

Add input validation and structured error responses for POST /v1/events.
Prevents panics on invalid payloads and improves observability.

Addresses security requirement from CODING_STANDARDS.md
```

## Pull Requests
- Reference relevant decision-log entries and update them if the change adds new assumptions.
- Provide tests and documentation for new features or behavior changes.
- Run the CI workflow locally (`make ci-check`) before opening the PR.
- Include screenshots or API examples for UI or endpoint changes.
- Link related issues and provide context for reviewers.

### PR Checklist
- [ ] Tests added/updated for new functionality
- [ ] Documentation updated (README, API docs, etc.)
- [ ] Security/privacy considerations addressed
- [ ] Performance impact assessed (if applicable)
- [ ] Decision log updated for significant design choices
- [ ] Compliance impact evaluated (if applicable)

## Decision Log Discipline
Every design or process choice with lasting impact must be captured in `decision-log.md`. Include: date, decision, rationale, and consequences.

Use this format:
```
| 2025-MM-DD | Brief decision summary | Why this choice was made | What it means for the project |
```

## Code Review Expectations
- Focus on correctness, determinism, and explainability.
- Verify that security/privacy considerations are addressed for input parsing and storage.
- Ensure documentation updates accompany code changes.
- Check for proper error handling and logging.
- Validate test coverage and quality.

## Development Setup

### First-time Setup
```bash
# Clone repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Install dependencies
go mod download
cd cbad-core && cargo fetch && cd ..

# Set up development environment
cp .env.example .env
$EDITOR .env

# Run initial build
make build
```

### Development Workflow
```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make changes and test locally
make test
make ci-check

# Commit with proper message format
git commit -m "feat(scope): brief description"

# Push and create PR
git push origin feature/your-feature-name
```

## Testing Requirements

### Unit Testing
- All new functions require unit tests
- Target ≥80% code coverage for core packages
- Use table-driven tests for multiple scenarios
- Include edge cases and error conditions

### Integration Testing
- Test end-to-end workflows
- Validate API contract compliance
- Test FFI boundary between Go and Rust
- Verify Docker composition works

### Performance Testing
- Benchmark critical paths
- Document performance baselines
- Test memory usage patterns
- Validate throughput targets

## Documentation Standards

### Code Documentation
- All public functions require clear docstrings
- Include usage examples for complex APIs
- Document mathematical formulas and algorithms
- Explain security/privacy considerations

### Architecture Documentation
- Use ASCII diagrams for broad compatibility
- Update ARCHITECTURE.md for significant changes
- Document configuration options thoroughly
- Include deployment and operational guidance

## Security Guidelines

### Input Validation
- Treat all inputs as untrusted
- Validate data lengths and formats
- Sanitize before processing or storage
- Use structured error responses

### Privacy Considerations
- Support configurable data redaction
- Document PII handling procedures
- Implement secure storage patterns
- Provide audit trail capabilities

### Dependency Management
- Pin dependency versions for reproducibility
- Regularly update for security patches
- Document security-critical dependencies
- Use trusted package sources only

## Compliance Considerations

When contributing features that affect compliance:
- Review relevant compliance templates (DORA, NIS2, Runtime AI)
- Update evidence bundle generation if needed
- Document audit trail implications
- Consider regulatory requirements

## Community Guidelines

### Communication
- Be respectful and constructive in discussions
- Ask questions if requirements are unclear
- Share knowledge and help other contributors
- Use GitHub issues for bug reports and feature requests

### Quality Standards
- Write clean, readable code
- Follow established coding standards
- Provide thorough testing
- Maintain excellent documentation

## Release Process

### Version Numbering
We follow semantic versioning (SemVer):
- MAJOR: incompatible API changes
- MINOR: backwards-compatible functionality
- PATCH: backwards-compatible bug fixes

### Release Checklist
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Performance benchmarks validated
- [ ] Security review completed
- [ ] Compliance impact assessed
- [ ] Migration guide provided (if needed)

## Getting Help

### Resources
- **Documentation**: Check `docs/` directory
- **Issues**: Search existing GitHub issues
- **Discussions**: Use GitHub Discussions for questions
- **Security**: Email security@driftlock.com for vulnerabilities

### Communication Channels
- GitHub Issues: Bug reports and feature requests
- GitHub Discussions: Questions and general discussion
- Pull Request Reviews: Technical feedback and collaboration

## Contributor License Agreement
By submitting a contribution, you agree that your work will be distributed under the Apache 2.0 License and that you have the right to grant this license.
