# Changelog

## [0.5.4](https://github.com/nu0ma/spemu/compare/v0.5.3...v0.5.4) (2025-07-05)

### Bug Fixes

* fix GitHub Actions workflow configuration ([dad8602](https://github.com/nu0ma/spemu/commit/dad86022c20f7d4f1c230eec88cf644d2b6345b5))

## [0.5.3](https://github.com/nu0ma/spemu/compare/v0.5.2...v0.5.3) (2025-07-05)

### Bug Fixes

* fix GitHub Actions workflow configuration ([d6f72b3](https://github.com/nu0ma/spemu/commit/d6f72b3245b65ef06394534c8502e2e166b7ee8b))

## [0.5.2](https://github.com/nu0ma/spemu/compare/v0.5.1...v0.5.2) (2025-07-05)

### Bug Fixes

* fix tag yaml GitHub Actions workflow ([54ee940](https://github.com/nu0ma/spemu/commit/54ee940accab7b3febc381739d87b28979d4d6c4))

## [0.5.1](https://github.com/nu0ma/spemu/compare/v0.5.0...v0.5.1) (2025-07-05)

### Features

* integrate GoReleaser for automated releases and distribution ([daea5f1](https://github.com/nu0ma/spemu/commit/daea5f127fa9d46f0756e4568d9c085b98eb3d46))
  - Add .goreleaser.yml configuration for cross-platform builds
  - Create GitHub Actions release workflow triggered by tags
  - Enable automated GitHub Releases with binaries and checksums
  - Add Homebrew formula generation for `brew install nu0ma/spemu/spemu`
  - Configure Docker image publishing to ghcr.io
  - Remove obsolete manual tag workflow and build targets
  - Simplify release process to git tag → automated distribution

## [0.5.0](https://github.com/nu0ma/spemu/compare/spemu-v0.4.1...v0.5.0) (2025-07-05)

### Features

* add manual tag creation workflow for flexible release options ([1972c33](https://github.com/nu0ma/spemu/commit/1972c33e7071330397f409c7f9e4c80b6ab32567))

## [0.4.1](https://github.com/nu0ma/spemu/compare/spemu-v0.4.0...spemu-v0.4.1) (2025-07-04)

### Bug Fixes

* replace Go setup script with shell script for integration tests ([588fc9f](https://github.com/nu0ma/spemu/commit/588fc9f))

### Documentation

* add CLAUDE.md for Claude Code integration ([ec9aa1f](https://github.com/nu0ma/spemu/commit/ec9aa1f))

## [0.4.0](https://github.com/nu0ma/spemu/compare/v0.3.0...spemu-v0.4.0) (2025-07-04)

### Features

* reduce code complexity and remove redundant files ([cd24997](https://github.com/nu0ma/spemu/commit/cd24997))

## [0.3.0](https://github.com/nu0ma/spemu/compare/v0.2.0...v0.3.0) (2025-07-04)

### Features

* add --version flag to display version information ([8aa3cd6](https://github.com/nu0ma/spemu/commit/8aa3cd66ba8a36e6a1bb3ad7930b491b54d3e1d8))
  - Add --version flag to main.go that displays current version
  - Update Makefile to inject version from version.txt during build
  - Update help text to include version option
  - Version is set via ldflags during build process

## [0.2.0](https://github.com/nu0ma/spemu/compare/spemu-v0.1.0...spemu-v0.2.0) (2025-07-04)

### Features

* add automated versioning and release system ([a76d200](https://github.com/nu0ma/spemu/commit/a76d200edef57b4735a1b546585c26eaead3c004))

### Bug Fixes

* improve GitHub Actions permissions for automated releases ([30f2162](https://github.com/nu0ma/spemu/commit/30f2162d4097006164216e867dd514e99618abfc))
* improve release-please workflow configuration ([58abf8e](https://github.com/nu0ma/spemu/commit/58abf8e63b6aa30cda99fd4e89c81979a7ea4fcf))