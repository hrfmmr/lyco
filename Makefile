PRJ=$(shell basename "$(PWD)")

.PHONY: generate
## generate: Execute go generate
generate:
	go generate ./...

.PHONY: build
## build: Build app
build: generate
	go build -o .build/app ./cmd/lyco

.PHONY: run
## run: Run app
run: build
	.build/app

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PRJ)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
