#!/usr/bin/make

.DEFAULT_GOAL := all

.PHONY: setup
setup:
	@go get golang.org/x/lint/golint
	@go get golang.org/x/tools/cmd/goimports
	@go get github.com/securego/gosec/v2

GOFILES=$(shell find . -type f -name '*.go' -not -path "./.git/*")

.PHONY: fmt
fmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -d -s -e $(GOFILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || (echo "gofmt failed:" | cat - $(FMT_LOG) && false)

.PHONY: imports
imports:
	$(eval IMP_LOG := $(shell mktemp -t goimp.XXXXX))
	@$(GOPATH)/bin/goimports -d -e -l $(GOFILES) > $(IMP_LOG) || true
	@[ ! -s "$(IMP_LOG)" ] || (echo "goimports failed:" | cat - $(IMP_LOG) && false)

.PHONY: lint
lint:
	@$(GOPATH)/bin/golint -set_exit_status $(shell go list ./...)

.PHONY: verify
verify:
	@make -s fmt
	@make -s imports
	@make -s lint

.PHONY: sec
sec:
	@gosec -quiet ./...

.PHONY: bin
bin:
	go build -o ./dist/route-service main.go

.PHONY: test
test:
	@go test ./... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: all
all:
	@make -s bin test verify sec

