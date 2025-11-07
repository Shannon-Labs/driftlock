# GitHub Release Preparation Checklist

This checklist ensures DriftLock is ready for a public GitHub release.

## Pre-Release Verification

### Code Quality
- [ ] All tests passing (`make test`)
- [ ] No compilation errors (`make build`)
- [ ] Linting passes (`make lint`)
- [ ] Security scan clean (`make security-scan`)
- [ ] No hardcoded secrets or credentials
- [ ] No proprietary code or dependencies

### Documentation
- [ ] README.md updated with current information
- [ ] CHANGELOG.md created with release notes
- [ ] CONTRIBUTING.md guidelines clear and accurate
- [ ] Installation guide (docs/installation.md) up to date
- [ ] API documentation (docs/API.md) complete
- [ ] Architecture documentation (docs/ARCHITECTURE.md) accurate
- [ ] Security policy (SECURITY.md) in place
- [ ] Code of conduct (CODE_OF_CONDUCT.md) present

### Configuration
- [ ] `.env.example` contains only necessary variables
- [ ] No hardcoded API keys or secrets in codebase
- [ ] Docker configuration works without external dependencies
- [ ] CI/CD workflows configured correctly
- [ ] GitHub issue templates created
- [ ] GitHub pull request template created

### Licensing
- [ ] Apache 2.0 license file present (LICENSE)
- [ ] Copyright headers on source files (if applicable)
- [ ] Third-party licenses documented
- [ ] No GPL or copyleft dependencies in core

### Repository Setup
- [ ] Repository visibility set to public
- [ ] Branch protection rules configured for `main`
- [ ] Required status checks enabled
- [ ] Code review requirements set
- [ ] Issue labels created and organized
- [ ] Repository description updated
- [ ] Topics/tags added to repository

### Release Preparation
- [ ] Version number updated in relevant files
- [ ] CHANGELOG.md updated with version and date
- [ ] Release notes prepared
- [ ] Docker images built and tested
- [ ] Installation instructions tested on clean environment
- [ ] Quick start guide verified to work

### Post-Release Tasks
- [ ] Create GitHub release with tag
- [ ] Upload release artifacts
- [ ] Announce on social media/channels
- [ ] Update project website (if applicable)
- [ ] Monitor for issues and respond promptly

## OSS Best Practices Compliance

### Repository Health
- [ ] Clear project description and purpose
- [ ] Up-to-date README with badges
- [ ] Clear contribution guidelines
- [ ] Issue templates for bugs and features
- [ ] Pull request template
- [ ] Code of conduct

### Technical Requirements
- [ ] No large binary files in repository
- [ ] `.gitignore` properly configured
- [ ] No sensitive data in git history
- [ ] Dependencies are up to date
- [ ] Known vulnerabilities addressed

### Community
- [ ] Responsive maintainers identified
- [ ] Clear governance model
- [ ] Roadmap or project vision shared
- [ ] Communication channels established

## Deployment Verification

### Docker Deployment
```bash
# Test clean build
git clone <repo-url> driftlock-test
cd driftlock-test
cp .env.example .env
# Edit .env with test configuration
docker compose up -d

# Verify services
 curl http://localhost:8080/healthz
# Should return 200 OK

# Test dashboard
curl http://localhost:3000
# Should return HTML
```

### Manual Build
```bash
# Test manual build process
make setup
make build
make test

# Verify binaries exist
ls -la bin/
```

## Final Checks

### Before Making Public
- [ ] Remove any internal comments or TODOs
- [ ] Verify no internal URLs or references
- [ ] Check for any proprietary code
- [ ] Ensure all tests pass in CI
- [ ] Verify documentation is accurate
- [ ] Test installation from scratch
- [ ] Have another team member review

### Release Day
- [ ] Create release tag (v0.1.0)
- [ ] Build and push Docker images
- [ ] Create GitHub release
- [ ] Post announcements
- [ ] Monitor for early user feedback

## Emergency Rollback Plan

If critical issues are found after release:

1. **Immediate**: Pin known-good version in documentation
2. **Within 1 hour**: Assess issue severity and impact
3. **Within 4 hours**: Prepare fix or rollback plan
4. **Within 24 hours**: Release patch or updated version
5. **Communicate**: Keep community informed of progress

---

**Note**: This checklist should be reviewed and updated for each release. Some items may not apply to all releases.
