PKG := ./cmd/lyco
BIN := lyco
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE=on

.PHONY: bootstrap
## bootstap: Set up for development
bootstrap:
	go install github.com/google/wire/cmd/wire@v0.4.0
	go install golang.org/x/tools/cmd/stringer@latest

.PHONY: generate
## generate: Execute go generate
generate:
	go generate ./...

.PHONY: build
## build: Build app
build: generate
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) $(PKG)

.PHONY: run
## run: Run app
run: build
	./$(BIN)

.PHONY: show-version
## show-version: Show current version of app
show-version: $(GOBIN)/gobump
	@gobump show -r .

$(GOBIN)/gobump:
	@cd && go get github.com/x-motemen/gobump/cmd/gobump

.PHONY: build-archive
## build-archive: Cross compile and archiving artifacts
build-archive: $(GOBIN)/goxz
	goxz -n $(BIN) -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) $(PKG)

$(GOBIN)/goxz:
	cd && go get github.com/Songmu/goxz/cmd/goxz

.PHONY: test
## test: Test
test: build
	go test -v ./...

.PHONY: clean
## clean: Clean build artifacts
clean:
	rm -rf $(BIN) goxz
	go clean

.PHONY: release
## release: Upload artifacts to GitHub Release
release: $(GOBIN)/ghr
	ghr "v$(VERSION)" goxz

$(GOBIN)/ghr:
	cd && go get github.com/tcnksm/ghr

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in lyco:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
