cwd := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
gobin := $(abspath $(cwd)/.gobin)

.PHONY: start
start:
	$(info starting chi server...)
	@go run cmd/server/main.go
