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
