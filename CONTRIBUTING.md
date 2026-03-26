# Contributing to gocontainerruntime

First off, thank you for considering contributing to `gocontainerruntime`! It's people like you that make the open-source community such an amazing place to learn, inspire, and create.

## Development Environment

To contribute to this project, you will need:

- **Go >= 1.22** (Standard Go compiler)
- **Linux Operating System** (Required for namespaces and cgroups)
- **Root access** (Required to run the container)
- **Make** (Optional but recommended for build automation)

## Code Style and Conventions

We follow the standard Go coding conventions:

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format your code
- Document all exported functions and types with high-quality GoDoc comments
- Keep functions and packages focused (Single Responsibility Principle)

## How to Contribute

1. **Fork the Repository** and create your branch from `main`.
2. **If you've added code** that should be tested, add tests.
3. **If you've changed APIs**, update the documentation.
4. **Ensure the test suite passes** by running `make test`.
5. **Make sure your code lints** by running `staticcheck ./...` (if available).
6. **Submit a Pull Request** describing your changes and the problem they solve.

## Commit Message Conventions

We use a simplified version of [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code changes that neither fix a bug nor add a feature
- `test:` for adding missing tests or correcting existing tests

## Need Help?

If you have any questions, feel free to open an issue or contact the maintainers.
