cwd := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
gobin := $(abspath $(cwd)/.gobin)

.PHONY: start
start:
	$(info starting chi server...)
	@go run cmd/server/main.go

.PHONY: gen
gen:
	$(info generating swagger...)
	@oapi-codegen -package api -generate types,strict-server,chi-server,client,spec -o gen/oapi.go spec/api.yaml
