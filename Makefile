IMG ?= hello-world-service:latest

build: 
	go build -o hello-world-service main.go

run:
	go run main.go

docker-build:
	docker build --tag ${IMG} .

docker-push:
	docker push ${IMG}

docker-run:
	docker run --publish 8080:8080 hello-world-service