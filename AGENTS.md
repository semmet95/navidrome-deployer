# Navidrome Deployer - Copilot PR Review Instructions

This document provides guidance for Copilot when reviewing pull requests targeting the `main` branch.

## General Review Guidelines

When reviewing PRs for the main branch, ensure:
- Code adheres to the project's architecture and patterns
- Changes maintain backward compatibility where applicable
- All modifications are properly documented
- Testing requirements are met

## Key Practices for Main Branch PRs

### 1. Helm Chart Version and Dependency Management

**Requirement:** Whenever the `navidrome-deployer` Helm chart is updated, the chart version **must** be incremented.

**What to check:**
- Changes to any files in `charts/navidrome-deployer/` (Chart.yml, values.yml, templates/*, etc.).
- Verify that `charts/navidrome-deployer/Chart.yml` has an updated version field.
- Use semantic versioning: increment patch version for bug fixes, minor for new features, major for breaking changes.
- Flag if chart changes are present but version remains unchanged.
- `appVersion` field in `charts/navidrome-deployer/Chart.yml` should be the same as the Navidrome image tag specified in `charts/navidrome-deployer/values.yaml`.
- `test/helmfile.yaml` and `helmfile.yaml` should have identical dependencies.
- `navidrome` chart version in `helmfile.yaml` should be the latest or next release version of `navidrome-deployer`.
- `filebrowser` version specified in `filebrowser.imageURI` field in `charts/navidrome-deployer/values.yaml` file should be the same as the `FILEBROWSER_VERSION` set in `Dockerfile.filebrowser` file.
- If `Dockerfile.filebrowser` file is updated `Release Packages` workflow should push the image with updated tag. `filebrowser.reconfigImageUri` field in `charts/navidrome-deployer/values.yaml` file should also be updated accordingly.

### 2. E2E Test Coverage for Features and Deployments

**Requirement:** Any new feature or deployment strategy must be covered by the e2e-tests workflow.

**What to check:**
- New deployment strategies or configuration options require corresponding e2e tests
- Feature branches should include test coverage additions in `.github/workflows/`
- Verify tests in the e2e-tests workflow validate the new behavior
- Flag missing e2e test coverage for new features or deployment changes

### 3. Mage for Environment Setup

**Requirement:** Use Mage (Go-based task runner) wherever possible for environment setup and minimize shell scripts.

**What to check:**
- New environment setup logic should be added to `magefiles/` rather than scripts
- Review `scripts/` directory changes carefully - new scripts should be justified
- Suggest migrating shell scripts to Mage tasks where appropriate
- Prefer Mage for build, test, and environment initialization tasks

### 4. Reproducible Environment Setup

**Requirement:** Any environment setup logic added to GitHub workflows must be easily reproducible locally.

**What to check:**
- GitHub Actions workflow setup steps should have corresponding Mage tasks or scripts
- Documentation should explain how to reproduce the workflow locally
- Setup commands should work identically in CI and local environments
- Avoid hardcoded paths, environment-specific configurations, or CI-only tooling when possible

## Review Checklist for Main Branch

- [ ] Chart version updated if `charts/navidrome-deployer/` modified
- [ ] `appVersion` and Navidrome image tag should be the same
- [ ] New features/deployments include e2e test coverage
- [ ] Environment setup uses Mage tasks instead of shell scripts where applicable
- [ ] Workflow setup is reproducible locally with clear documentation
- [ ] Code quality and testing standards are maintained
- [ ] Documentation is updated for user-facing changes

## Questions to Ask During Review

1. **For chart changes:** "Was the chart version updated in Chart.yml?"
2. **For new features:** "Are there corresponding e2e tests validating this change?"
3. **For setup changes:** "Could this be implemented as a Mage task instead?"
4. **For CI/CD changes:** "Can I run these setup steps locally for testing?"

---

*Last updated: December 2025*
