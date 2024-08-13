# Overview

This is used to generate an artifacthub compatible changelog between two tags of a repository.
For example:

```
annotations:
  artifacthub.io/changes: |
    - kind: added
      description: Support for Tekton pipelines
      links:
          - name: Github Issue
            url: https://github.com/artifacthub/hub/issues/1485
    - kind: added
      description: Versions index to changelog modal
    - kind: added
      description: Allow publishers to include screenshots in packages
```

This needs to be added to a helm chart release on artifacthub to generate changelogs between released versions.

You can read further about this here: https://blog.artifacthub.io/blog/changelogs/

# Usage

```
 go run main.go \
  --to v1.7.0 \
  --from v1.0.0 \
  --repoURL helm-artifacthub-chglog/fixtures/go-siva \
  --output helm-artifacthub-chglog/_example-outputs/go-siva.yaml
```

Will generate: `helm-artifacthub-chglog/_example-outputs/go-siva.yaml` and use the default resolver for conventional commits.

```
 go run main.go \
  --to router@0.13.0 \
  --from router@0.89.2 \
  --repourl helm-artifacthub-chglog/fixtures/cosmo \
  --paths router \
  --output helm-artifacthub-chglog/_example-outputs/cosmo-router@0.89.2-router@0.13.0.yaml
```

will generate: `helm-artifacthub-chglog/_example-outputs/cosmo-router@0.89.2-router@0.13.0.yaml`.

those commits are all those between two tags `from` `to` and related to paths in the repository.

All commands can be run with [task](https://taskfile.dev/).

e.g.:

```
# Quick Start
task # will run all tests and build the binary


# Fixtures tests
task test:fixture:conventional-changelog # runs the fixture tests on https://github.com/conventional-changelog/conventional-changelog.git

# E2E tests and build
task test:e2e build # run the tests in the e2e folder and build the binary
```
