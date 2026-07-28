VERSION = $(shell date '+%Y.%-m.%-d')

.DEFAULT_GOAL := build
.PHONY: build clean lint release update

clean:
	@rm -rf dist/ cmd/static/*.css cmd/static/*.js

update:
	@go get -u ./cmd
	@go mod tidy

build:
	@goreleaser build --single-target --snapshot --clean

release:
	@git tag -f v$(VERSION)
	@goreleaser release --clean

lint:
	@go vet ./cmd
	@-golangci-lint run
	@gofmt -d ./cmd
	@biome check static/*.js
	@prettier --check static/*.scss
