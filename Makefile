.DEFAULT_GOAL: help
SHELL := /bin/bash

PROJECTNAME := "pears-dht"
KIND_CLUSTER_NAME := "pears"
WAIT_FOR_KIND_READY = '{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'

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

## create-cluster: Create a kind cluster named "pears"
create-cluster: check-kind
ifeq (1, $(shell kind get clusters | grep ${KIND_CLUSTER_NAME} | wc -l))
	@echo "Cluster already exists - deleting it to start from clean cluster"
	kind delete cluster --name ${KIND_CLUSTER_NAME}
endif
	@echo "Creating Cluster"
	kind create cluster --name ${KIND_CLUSTER_NAME} --image=kindest/node:v1.24.0
	until kubectl get nodes -o jsonpath="${WAIT_FOR_KIND_READY}" 2>&1 | grep -q "Ready=True"; do sleep 5; echo "--------> waiting for cluster node to be available"; done

## run-server-mode: Create a pears node as the server
run-server-mode:

## build: Build the container image
build: check-docker
	@docker build --no-cache --pull -f build/packages/Dockerfile -t ${PROJECTNAME}:local-build .

## dev-setup: Create a complete dev-setup with k8s cluster running 4 pears dht
dev-setup: build create-cluster run-server-mode
