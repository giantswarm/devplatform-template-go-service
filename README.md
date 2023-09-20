# Go web service template

![validation and test result](https://github.com/giantswarm/devplatform-template-go-service/actions/workflows/validate-test.yaml/badge.svg?event=push)

## NOTE

This repo is currently a private repo and a self-sufficient one. In other words, it's currently only an internal preview
of a repo that we want to make public, but also turn into a real template, that can be used to create other projects of
the same type using a dev platform workflow.

## Intro

This repo is meant as a template base for golang based web service projects, that are deployed to a Kubernetes cluster
using a Helm Chart. The repo follows the rule of fast local builds and developer feedback: tools configured for the CI
process in the repo are also installable on local dev machines, allowing for rapid feedback loops, without waiting for
the CI.

## Features included

- automatically build go binaries, a container image and a Helm chart
- upload build artifacts to GitHub: create a release for the binaries, upload the container image and the Helm chart
  to GitHub's OCI registry
- included security: vulnerability scans for go sources, generation of SBoM, singing artifacts with `cosign`
- included automated dependency updates based on [renovate](renovatebot.com)
- included linting and validation for multiple types of artifacts, including golang, markdown, Kubernetes objects, ...

### How it works

#### Running validation and linting tools

This repo uses [pre-commit](https://pre-commit.com/) to run a variety of static code quality analysis tools, both in
CI and local runs. Please [consult the config file](.pre-commit-config.yaml) to check what is run by default.

The checks configured include:

- general good practices for code files (ie. consistent line endings)
- `golang` validation based on the [golangci-lint](https://golangci-lint.run/) tool
- Dockerfile, markdown, JSON, YAML and shell scripts linting
- [gitleaks](https://github.com/gitleaks/gitleaks) tool to protect developers from committing secrets to the repo
- [helm-docs](https://github.com/norwoodj/helm-docs) to automatically generate Helm chart's README.md file with
  chart's configuration options generated from `values.yaml` file

#### Building binaries and container images

For this purpose, the repo uses [goreleaser](https://goreleaser.com/). Please consult [its config file](.goreleaser.yaml)
to check the included configuration and tune it to your needs.

By default, the `goreleaser` configuration includes:

- building the go binary for 3 CPU architectures: `amd64`, `arm` and `arm64`
- building multi-architecture container image with the binary built in the previous step
- generating SBoM manifest
- signing all the artifacts with `cosign`
- creating a GitHub release with automated changelog based on git commits

#### Building a Helm chart

To run a set of code linting and validation tools, the setup uses
[app-build-suite](https://github.com/giantswarm/app-build-suite/). The build process includes

- [ct](https://github.com/helm/chart-testing) Helm chart testing tool
- Kubernetes objects validation with [kube-linter](https://github.com/stackrox/kube-linter)
- actual Helm chart packaging

### GitHub actions

In the repo, there are already GitHub actions pre-configured. Their purpose is as follows:

- [.github/workflows/release.yaml](.github/workflows/release.yaml): triggered on a tag push and does full project release
- [.github/workflows/validate-test.yaml](.github/workflows/validate-test.yaml): triggered on each commit, runs basic
   code validation
- [.github/workflows/update-pre-commit.yaml](.github/workflows/update-pre-commit.yaml): ran on schedule to detect
  automatic updates to `pre-commit` actions

Additionally, renovate is [configured](renovate.json) to automatically handle dependencies updates.

## Setting up local development environment

### Minimal setup

The `pre-commit` tools is mandatory. It will guard you from creating inconsistencies between you and other team members,
lower CI costs by limiting the necessary number of runs and even do some basic security checks.

To install `pre-commit` and its dependencies, please refer the following docs:

- <https://pre-commit.com/#install>
- <https://golangci-lint.run/usage/install/>
- <https://github.com/norwoodj/helm-docs#installation>
- obviously, a working `golang` installation

Once you have the dependencies installed, you have to register `pre-commit` as your local hooks for the git repo
(you need to do this only once):

```bash
pre-commit install --install-hooks
```

To run all the `pre-commit` checks without creating a commit, run:

```bash
pre-commit run -a
```

### Setup for full build workflow

By installing the same set of tools GitHub actions run, you can have a local dev environment that works the same
as CI/CD and provides you with faster feedback loop.

Install additional tools:

- [docker](https://www.docker.com/products/docker-desktop/)
- [goreleaser](https://goreleaser.com/)
- [syft](https://github.com/anchore/syft)
- [cosign](https://github.com/sigstore/cosign)
- [app-build-suite](https://github.com/giantswarm/app-build-suite/releases) - download the latest script (dockerized)

Now, to build the go binary and docker image locally, run:

- `goreleaser release --verbose --snapshot`

To build the Helm chart, run:

- `dabs.sh -c helm`

## Creating a release

To release a new version of an app (create a GitHub release, container image and Helm chart), create and push a tag
to trigger the process, ie.

```bash
git tag v0.1.0
git push origin v0.1.0
```
