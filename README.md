# hello-world-service

A simple HTTP server that serves "Hello world" on http://localhost:8080/api/hello-world-service/v0/ping

**To run without docker:**
1. Install Golang 1.15.
2. Run `make run`

**To run with docker:**
1. Start the docker daemon
2. Export the following environment variables (needed for pulling the base image from redhat.registry.io):
   1. RH_REGISTRY_USER (redhat.registry.io service account user)
   2. RH_REGISTRY_TOKEN (redhat.registry.io service account token)
3. Run `make docker-build`
4. Run `make docker-run`

**Testing and Linting:**

*Unit tests:*
```
make test
```
*Lint*:
```
make lint
```

**To push out to a quay.io repo:**
1. Export the following environment variables
   1. QUAY_USER (quay.io user name)
   2. QUAY_TOKEN (quay.io token/ encrypted CLI password)
   3. IMG="quay.io/${QUAY_USER}/hello-world-service:latest"
2. Run `make docker-push`
