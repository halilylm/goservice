.PHONY: run

VERSION := 1.0
KIND_CLUSTER := service-cluster
VCS_REF := $(shell git rev-parse HEAD)

run-products-api:
	go run app/services/products-api/main.go

build-products-api:
	go build -ldflags "-X main.build=local" -o products-api ./app/services/products-api

build-products-api-docker:
	docker build --no-cache \
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

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-apply:
	cat zarf/k8s/base/service-pod/base-service.yaml | kubectl apply -f -

kind-load:
	kind load docker-image products-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-restart:
	kubectl rollout restart deployment service-pod --namespace=service-system

kind-update-products-api: build-products-api-docker kind-restart

kind-logs:
	kubectl logs -l app=service-system --all-containers=true -f --tail=100