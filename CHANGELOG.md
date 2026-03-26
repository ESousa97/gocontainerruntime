# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete professional documentation suite (README, CONTRIBUTING, LICENSE, etc.).
- `Makefile` for standardizing build and run processes.
- Enhanced GoDoc comments in `main.go`.
- Technical architecture documentation in `docs/architecture.md`.

## [1.0.0] - 2026-03-26

### Added
- Initial implementation of the Go container runtime.
- Support for Linux Namespaces (UTS, PID, Mount, Network).
- Basic Cgroups support for Memory and CPU limits.
- Alpine Linux rootfs download and extraction logic.
- Command-line interface using Cobra.
- Chroot-based filesystem isolation.
