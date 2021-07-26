IMG ?= hello-world-service:latest

build: 
	go build -o hello-world-service main.go

run:
	go run main.go

# Docker build, push and run
docker-build:
# Base image for go is pulled from registry.redhat.io
	docker login -u="${RH_REGISTRY_USER}" -p="${RH_REGISTRY_TOKEN}" registry.redhat.io
	docker build --tag ${IMG} .

docker-push:
	docker login -u="${QUAY_USER}" -p="${QUAY_TOKEN}" quay.io
	echo ${IMG}
	$(MAKE) docker-build
	docker push ${IMG}

docker-run:
	docker run --publish 8080:8080 hello-world-service

test:
	go test

vet:
	go vet ./...

staticcheck: 
	staticcheck ./...

lint: vet staticcheck
# Note: The Golint linter is deprecated and frozen. As per the docs (https://github.com/golang/lint) there's no drop-in replacement for it, but tools such as Staticcheck and go vet should be used instead.

# Open API 3.0 Spec
generate-docs:
	go run cmd/spec/main.go