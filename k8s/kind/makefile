SHELL := /bin/bash

all:
	(cd ../.. && make all)

getlive-api:
	(cd ../.. && make getlive-api)

metrics:
	(cd ../.. && make metrics)

test:
	(cd ../.. && make test)

up:
	kind create cluster --name ardan-starter-cluster --config kind-config.yaml

down:
	kind delete cluster --name ardan-starter-cluster

load:
	kind load docker-image gcr.io/getlive-backend-api/getlive-api-amd64:1.0 --name ardan-starter-cluster
	kind load docker-image gcr.io/getlive-backend-api/metrics-amd64:1.0 --name ardan-starter-cluster
	# kind load docker-image openzipkin/zipkin:2.11 --name ardan-starter-cluster
	# kind load docker-image postgres:11.1-alpine --name ardan-starter-cluster

services:
	kubectl create -f deploy-postgres.yaml
	kubectl create -f deploy-getlive-api.yaml
	@echo ======================================================================

update-getlive-api:
	# Build a new version using 1.1
	kind load docker-image gcr.io/getlive-backend-api/getlive-api-amd64:1.1 --name ardan-starter-cluster
	kubectl set image pod <POD_NAME> getlive-api=gcr.io/getlive-backend-api/getlive-api-amd64:1.1
	kubectl delete pod <POD_NAME>

get-pods:
	kubectl get pods

logs:
	kubectl logs <POD_NAME> getlive-api

status:
	kubectl get nodes
	kubectl get pods
	kubectl get services getlive-api
	kubectl logs <POD_ID> -c getlive-api
	@echo ======================================================================

shell:
	# kubectl get pods
	# kubectl exec -it <POD NAME> --container getlive-api  -- /bin/sh
	# ./admin --db-disable-tls=1 migrate
	# ./admin --db-disable-tls=1 seed
	# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token
	# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
	# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users
	@echo ======================================================================

delete:
	kubectl delete services getlive-api
	kubectl delete deployment getlive-api	
	@echo ======================================================================