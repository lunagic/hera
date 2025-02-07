# hera

Development Environment Overseer

## Installation

```shell
go install github.com/lunagic/hera@latest
```

## Command Line Usage

```shell
hera
```

## Example Configuration File

```yaml
# .config/hera.yaml
services:
  backend:
    command: go run .
    watch:
      - Makefile
      - src/backend
      - go.mod
      - main.go
    exclude:
      - src/backend/dist
  frontend:
    command: npm run build
    watch:
      - Makefile
      - src/frontend
      - package.json
```
