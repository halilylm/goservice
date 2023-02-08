.PHONY: run

VERSION := 1.0
KIND_CLUSTER := service-cluster
VCS_REF := $(shell git rev-parse HEAD)

run-products-api:
	go run app/services/products-api/main.go

build-products-api:
	go build -ldflags "-X main.build=local" -o products-api ./app/services/products-api

docker-products-api:
	docker build \
	-f zarf/docker/dockerfile.products-api \
	-t service-api-amd64:$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg VCS_REF=$(VCS_REF) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.

kind-up:
	kind create cluster \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=service-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image service-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch

kind-status-postgres:
	kubectl get pods -o wide --watch --namespace=database-system

kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	#kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/service-pod | kubectl apply -f -

kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment service-pod

kind-status-service:
	kubectl get pods -o wide --watch

kind-update: docker-products-api kind-load kind-restart

kind-update-apply: docker-products-api kind-load kind-apply

kind-start-over: kind-down kind-up docker-products-api kind-load kind-apply

test:
	go test ./... -count=1
	staticcheck -checks=all ./...