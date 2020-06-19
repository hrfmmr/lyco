PRJ=$(shell basename "$(PWD)")
BIN := lyco

.PHONY: generate
## generate: Execute go generate
generate:
	go generate ./...

.PHONY: build
## build: Build app
build: generate
	go build -o .build/$(BIN) ./cmd/lyco

.PHONY: run
## run: Run app
run: build
	.build/$(BIN)

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PRJ)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
