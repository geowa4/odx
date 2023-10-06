# odx

Run one-off scripts and binaries

```shell
❯ cat ~/.config/odx.yaml
---
sources:
  utils:
    github: geowa4/ops-sop
    branch: patch-1
    path: v4/utils
aliases:
  prune:
    - utils
    - pruningfix.sh

❯ export ODX_GITHUB_TOKEN=... # personal token with all repo permissions

❯ go run main.go utils pruningfix.sh
Finding affected nodes
...

❯ go run main.go prune
Finding affected nodes
...
```