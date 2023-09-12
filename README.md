# Go web service template

## DRAFT

## Setting up local development environment

This way you can have local dev environment same as CI/CD.

Install git checks:

- pre-commit
- golangci-lint

If you want to build locally the same way as CI/CD:

- docker
- goreleaser
- syft
- cosign

Configure (first time only):

- `pre-commit install --install-hooks`

To build locally:

- `goreleaser release --verbose --snapshot`
