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
  --repoURL helm-artifacthub-chglog/fixtures/cosmo \
  --paths router \
  --output helm-artifacthub-chglog/_example-outputs/cosmo-router@0.89.2-router@0.13.0.yaml
```

Will generate: `helm-artifacthub-chglog/_example-outputs/cosmo-router@0.89.2-router@0.13.0.yaml`.

Those commits are all those between two tags `from` `to` and related to paths in the repository.
