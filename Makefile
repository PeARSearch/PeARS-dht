.DEFAULT_GOAL: help
SHELL := /bin/bash

PROJECTNAME := "pears-dht"

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo note: call scripts from /scripts

check-%: # detection of required software.
	@which ${*} > /dev/null || (echo '*** Please install `${*}` ***' && exit 1)

## build: Build the container image
build: check-docker
	@docker build --no-cache --pull -f build/packages/Dockerfile -t ${PROJECTNAME}:local-build .
