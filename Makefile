PRJ=$(shell basename "$(PWD)")

.PHONY: build
## build: Build app
build:
	go build -o .build/app .

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
