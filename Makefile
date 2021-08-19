# Copyright Red Hat

-include /opt/build-harness/Makefile.prow

S := $(shell uname)
UNAME_S := $(shell uname -s)
OS_SED :=
ifeq ($(UNAME_S),Darwin)
	OS_SED += ""
endif

IMAGE="quay.io/cloudservices/idp-configs-api"
IMAGE_TAG="latest"

KUBECTL=kubectl
NAMESPACE=default

check: check-copyright

check-copyright:
	@build/check-copyright.sh

build: 
	go build -o idp-configs-api main.go

run:
	go run main.go

# Docker build, push and run
docker-build:
# Base image for go is pulled from registry.redhat.io
	docker login -u="${RH_REGISTRY_USER}" -p="${RH_REGISTRY_TOKEN}" registry.redhat.io
	docker build --tag ${IMAGE}:${IMAGE_TAG} .

docker-run:
	docker run --publish 3000:3000 ${IMAGE}:${IMAGE_TAG}

test:
	go test -v -coverprofile=coverage.out

vet:
	go vet ./...

staticcheck: 
	staticcheck ./...

lint: vet staticcheck
# Note: The Golint linter is deprecated and frozen. As per the docs (https://github.com/golang/lint) there's no drop-in replacement for it, but tools such as Staticcheck and go vet should be used instead.

# OpenAPI 3.0 Spec
generate-docs:
	go run cmd/spec/main.go

bonfire-config-local:
	@cp default_config.yaml.local.example config.yaml
	@sed -i ${OS_SED} 's|REPO|$(PWD)|g' config.yaml

bonfire-config-github:
	@cp default_config.yaml.github.example config.yaml	

create-ns:
	$(KUBECTL) create ns $(NAMESPACE)

deploy-env:
	bonfire deploy-env -n $(NAMESPACE)

deploy-app:
	bonfire deploy idp-configs -n $(NAMESPACE)

scale-down:
	$(KUBECTL) scale --replicas=0 deployment/idp-configs-api-service -n $(NAMESPACE)

scale-up:
	$(KUBECTL) scale --replicas=1 deployment/idp-configs-api-service -n $(NAMESPACE)

restart-app:
	$(MAKE) scale-down NAMESPACE=$(NAMESPACE)
	sleep 5
	$(MAKE) scale-up NAMESPACE=$(NAMESPACE)	
