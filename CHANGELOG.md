# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Added smoke test for local provider functionality in `test/smoke_test.sh`
  - Verifies creation of a local file with expected content using OpenTofu
  - Includes proper cleanup of test artifacts
  - Uses the `tofu` CLI for OpenTofu operations

## [0.1.0] - 2025-8-15

### Added
- Initial project setup
