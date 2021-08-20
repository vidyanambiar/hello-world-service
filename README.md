[comment]: # ( Copyright Red Hat )
# idp-configs-api

Service to store and retrieve Identity Provider configurations. 

---
### Basic setup:

1. Install Golang 1.16
2. Run `make run-migrate`
3. Run `make run`
4. Access the API at http://localhost:3000/api/idp-configs-api/v0/ping

---

### Setup with Kubernetes (Minikube)

This setup utilizes the following tools. Follow the steps in the links for installation:
- [minikube](https://minikube.sigs.k8s.io/docs/)
- [Clowder](https://github.com/RedHatInsights/clowder)
- [Bonfire](https://github.com/RedHatInsights/bonfire)

Following the information above you should have Docker, a minikube cluster running with Clowder installed, and a Python environment with `bonfire` installed. Now move on to running the `idp-configs-api` application.

1. Clone the project.
```
git clone git@github.com:identitatem/idp-configs-api.git
```
2. Change directories to the project.
```
cd idp-configs-api
```
3. Setup your Python virtual environment.
```
pipenv install --dev
```
4. Enter the Python virtual environment to enable access to Bonfire.
```
pipenv shell
```
5. Setup access to the Docker enviroment within minikube, so you can build images directly to the cluster's registry.
```
eval $(minikube -p minikube docker-env)
```
6. Build the container image.
```
make docker-build
```
7. Create Bonfire configuration. To deploy from your local repository run the following:
```
make bonfire-config-local
```
The above command will create a file named `config.yaml` pointing to your local repository. At times you may need to update the branch which is referred to with the `ref` parameter (defaults to `main`).

Bonfire can also deploy from GitHub. Running the following command will setup the GitHub based configuration:
```
make bonfire-config-github
```
8. Setup test namespace for deployment.
```
make create-ns NAMESPACE=test
```
9. Deploy a Clowder environment (*ClowdEnviroment*) to the namespace with bonfire.
```
make deploy-env NAMESPACE=test
```
10. Deploy the application to the namespace.
```
make deploy-app NAMESPACE=test
```

Now the application should be running. You can test this by port-forwarding the app in one terminal and running a curl command in another as follows:

**Terminal 1**
```
kubectl -n test port-forward service/idp-configs-api-service 8000:8000
```
**Terminal 2**
```
curl -v http://localhost:8000/
```

You should get a 200 response.

---
### Testing and Linting:

Unit tests:
```
make test
```
Lint:
```
go get -u honnef.co/go/tools/cmd/staticcheck@latest
make lint
```

---
#### Viewing the OpenAPI 3.0 spec:

* Run `make generate-docs`

The API will serve the docs under a `/docs` endpoint.