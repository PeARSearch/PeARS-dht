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

## build: Build the container image
build: check-docker
	@docker build --no-cache --pull -f build/packages/Dockerfile -t ${PROJECTNAME}:local-build .

## run-pear: Run a local pear with docker in server mode(as the first node)
run-pear: build delete-docker-ps
	@docker run --name pear-1 -p 8080:8080 -p 8888:8888 -it ${PROJECTNAME}:local-build

## delete-pear: Delete the pear
delete-pear: delete-docker-ps delete-docker-image

## delete-docker-ps: Delete the pear process
delete-docker-ps:
	@docker rm pear-1 || echo "No pear container running"

## delete-docker-image:
delete-docker-image:
	@docker rmi ${PROJECTNAME}:local-build || echo "Image already cleaned up"

# ------ Kubernetes setup targets ----- #

## create-cluster: Create a kind cluster named "pears"
create-cluster: check-kind
ifeq (1, $(shell kind get clusters | grep ${KIND_CLUSTER_NAME} | wc -l))
	@echo "Cluster already exists - deleting it to start from clean cluster"
	kind delete cluster --name ${KIND_CLUSTER_NAME}
endif
	@echo "Creating Cluster"
	kind create cluster --name ${KIND_CLUSTER_NAME} --image=kindest/node:v1.24.0
	until kubectl get nodes -o jsonpath="${WAIT_FOR_KIND_READY}" 2>&1 | grep -q "Ready=True"; do sleep 5; echo "--------> waiting for cluster node to be available"; done

## load-docker-image: Load docker image to the kind cluster
load-docker-image:
	@kind load docker-image --name pears ${PROJECTNAME}:local-build

## run-server-pod: Create a pears node as the server
run-server-pod:
	@kubectl run pears --image=${PROJECTNAME}:local-build  --port=8080 --port=8888 --expose

## dev-setup: Create a complete dev-setup with k8s cluster running 4 pears dht
dev-setup: build create-cluster load-docker-image run-server-pod

## delete-kind-cluster:
delete-kind-cluster:
	@kind delete cluster --name ${KIND_CLUSTER_NAME} || echo "Cluster already deleted"

## k8s-cleanup: Cleanup all the created resources
k8s-cleanup: delete-docker-image delete-kind-cluster
