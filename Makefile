APP_NAME := montty
CURRENT_PLATFORM := $(shell printf '%s-%s' $$(go env GOOS GOARCH))
VERSION := $(shell date '+%Y.%-m.%-d')
GIT_COMMIT := $(shell git rev-parse --short HEAD)
PLATFORMS := $(sort darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-arm $(CURRENT_PLATFORM))
ASSETS := $(patsubst static/%.scss,cmd/static/%.css,$(wildcard static/*.scss))
ASSETS += $(patsubst static/%.js,cmd/static/%.js,$(wildcard static/*.js))

.DEFAULT_GOAL := build
.PHONY: clean update assets build build-all release lint $(PLATFORMS)

clean:
	@rm -rf dist/ $(ASSETS)

update:
	@go get -u ./cmd
	@go mod tidy

cmd/static/%.css: static/%.scss
	@sass --style=compressed --no-source-map $< $@

cmd/static/%.js: static/%.js
	@esbuild --bundle --minify --platform=browser --legal-comments=none --outfile=$@ $<

assets: $(ASSETS)

linux-arm: export GOARM=5
$(PLATFORMS): OUTPUT=$(APP_NAME)-$@$(if $(findstring windows,$@),.exe,)
$(PLATFORMS): export GOOS=$(word 1,$(subst -, ,$@))
$(PLATFORMS): export GOARCH=$(word 2,$(subst -, ,$@))
$(PLATFORMS):
	@go build \
		-C cmd \
		-trimpath \
		-buildvcs=false \
		-ldflags '-s -w -X main.version=$(VERSION)+$(GIT_COMMIT)' \
		-o '../dist/$(OUTPUT)'
	@echo $(OUTPUT)

build: assets $(CURRENT_PLATFORM)

build-all: MAKEFLAGS+=-j
build-all: $(PLATFORMS)

release: clean assets lint build-all
	@find dist -type f  | parallel 'xz -z9v {}'
	@rhash -r --printf '%{sha-256}  %f\n' dist > dist/SHA256SUMS
	@git tag -f 'v$(VERSION)'

lint:
	@go vet ./cmd
	@-golangci-lint run
	@gofmt -d ./cmd
