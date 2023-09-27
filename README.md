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
