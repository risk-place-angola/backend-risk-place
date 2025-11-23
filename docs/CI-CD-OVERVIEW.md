# CI/CD Overview

## 🎯 Purpose

This document provides a high-level overview of the CI/CD process for contributors. Detailed deployment configuration is maintained separately by the project maintainers.

## 🔄 Continuous Integration (CI)

### Trigger
Pull Requests to the `main` branch automatically trigger the CI workflow.

### What Gets Checked

1. **Security Scanning**
   - Uses `govulncheck` to scan for known vulnerabilities
   - Blocks merge if critical vulnerabilities are found

2. **Code Linting**
   - Runs `golangci-lint` with project configuration
   - Ensures code style consistency
   - Catches common mistakes and anti-patterns

3. **Unit Tests**
   - Executes all test suites
   - Generates coverage reports
   - Requires passing tests for merge

4. **Build Verification**
   - Compiles the application
   - Ensures no build-breaking changes

### Workflow File

The CI workflow is defined in `.github/workflows/ci.yaml`

### Running Locally

Before creating a PR, run these checks locally:

```bash
# Security scan
make sec-scan

# Lint
make lint

# Tests
make test

# Build
make build
```

## 🚀 Continuous Deployment (CD)

### Overview

When code is merged to `main`, an automated deployment process:

1. Builds a production Docker image
2. Pushes to a container registry
3. Deploys to production servers
4. Performs health checks
5. Rolls back if deployment fails

### Docker Image

The production image is built using:
- `Dockerfile.prod` - Multi-stage optimized build
- Alpine Linux base for minimal size
- Non-root user for security
- Health checks built-in

### For Contributors

**You don't need to worry about deployment!** Just:

1. Write code
2. Write tests
3. Create PR to `develop`
4. Address review feedback
5. Once merged to `main`, maintainers handle deployment

## 🏗️ Architecture

### Development Flow

```
Feature Branch → PR to develop → Code Review → Merge to develop
                                                      ↓
                                              Testing & Validation
                                                      ↓
                                              PR to main → CI
                                                      ↓
                                            Merge → CD (Maintainers)
                                                      ↓
                                                  Production
```

### Docker Configuration

**Development** (`docker-compose.yml`):
- PostgreSQL
- Redis
- Mailpit (email testing)
- Hot reload for development

**Production** (Maintainer-configured):
- Optimized builds
- Production databases
- Load balancing
- Health monitoring

## 🧪 Testing

### Local Testing

```bash
# Run all tests
make test

# Run specific test
go test -run TestUserValidation ./internal/domain/model

# With coverage
go test -cover ./...
```

### Docker Testing

```bash
# Build production image locally
docker build -f Dockerfile.prod -t backend-core:test .

# Run locally
docker run --rm -p 8090:8090 --env-file .env.dev backend-core:test

# Test health
curl http://localhost:8090/health
```

## 📊 Build Process

### Multi-Stage Build

1. **Builder Stage**:
   - Uses official Go image
   - Downloads dependencies
   - Compiles binary with optimizations
   - Strips debug symbols

2. **Final Stage**:
   - Minimal Alpine Linux
   - Copies only the binary
   - Runs as non-root user
   - ~20-30MB final image

### Build Optimization

- Layer caching for faster builds
- Only necessary files included
- No development dependencies
- Optimized binary size

## 🔒 Security

### Image Security

- ✅ No root user
- ✅ Minimal attack surface (Alpine)
- ✅ No unnecessary packages
- ✅ Regular base image updates
- ✅ Vulnerability scanning

### Code Security

- ✅ Dependency scanning
- ✅ Static analysis
- ✅ Secret detection (pre-commit hooks)
- ✅ HTTPS only in production
- ✅ Environment variable isolation

## 🤝 Contributing

### Before Creating a PR

1. **Ensure tests pass**:
   ```bash
   make test
   ```

2. **Run linter**:
   ```bash
   make lint
   ```

3. **Check for vulnerabilities**:
   ```bash
   make sec-scan
   ```

4. **Test build**:
   ```bash
   make build
   ```

### PR Guidelines

- Create PRs against `develop` branch
- Write descriptive commit messages
- Include tests for new features
- Update documentation as needed
- Reference related issues

### Review Process

1. Automated CI checks must pass
2. Code review by maintainers
3. Address feedback
4. Approval required before merge

## 📚 Additional Resources

### For Contributors

- [Contributing Guide](../CONTRIBUTING.md)
- [Code of Conduct](../CODE_OF_CONDUCT.md)
- [Development Setup](../README.md#getting-started)

### For Developers

- [API Documentation](../api/)
- [Architecture Diagrams](../docs/diagram/)
- [Domain Documentation](../docs/)

## ❓ FAQ

### Do I need Docker Hub access?

No, contributors don't need Docker Hub access. Images are built and published by CI/CD.

### How do I test deployment locally?

Use `docker-compose.yml` for local development. Production deployment is handled by maintainers.

### Where are production secrets stored?

Production configuration is maintained separately by project maintainers for security.

### Can I see deployment logs?

Public CI logs are visible in GitHub Actions. Production deployment logs are private.

### How long does CI take?

Typically 3-5 minutes for:
- Security scan: ~30s
- Lint: ~30s
- Tests: ~1-2min
- Build: ~1min

---

**Questions?** Open an issue or contact the maintainers.

**Last Updated**: November 22, 2025
