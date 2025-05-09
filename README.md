# hera

Development Environment Overseer

[![Validate](https://github.com/lunagic/hera/actions/workflows/validate.yml/badge.svg)](https://github.com/lunagic/hera/actions/workflows/validate.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/lunagic/hera.svg)](https://pkg.go.dev/github.com/lunagic/hera)
[![Go Report Card](https://goreportcard.com/badge/github.com/lunagic/hera)](https://goreportcard.com/report/github.com/lunagic/hera)

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
