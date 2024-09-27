GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

vault: fmt vault-build vault-start

vault-build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/vault-plugin-vcluster-operator cmd/vault-plugin-vcluster-operator/main.go

vault-start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

vault-enable:
	vault secrets enable -path=vcluster-operator vault-plugin-vcluster-operator

vault-clean:
	rm -f ./vault/plugins/vault-plugin-vcluster-operator

fmt:
	go fmt $$(go list ./...)

.PHONY: vault-build vault-clean fmt vault-start vault-enable


all: docker deploy

build:
	./ci/set-version.sh
	./ci/go-build.sh

docker: ##  Builds the application for amd64 and arm64
	docker buildx build --platform linux/amd64,linux/arm64 -t mheers/vcluster-operator:latest --push .

deploy: deploy-clean deploy-image install

deploy-clean: uninstall
	sleep 30

deploy-image:
	k3d image import mheers/vcluster-operator

install:
	go run main.go install --admin-username root --admin-password admin --image-pull-policy Never

uninstall:
	go run main.go uninstall

watch:
	kubectl logs -f deployments/vcluster-operator

forward:
	kubectl port-forward deployments/vcluster-operator -n vcluster-operator 8080:8080

server:
	export VCLUSTER_OPERATOR_K8S_INCLUSTER=false
	export VCLUSTER_OPERATOR_PORT=8080
	export VCLUSTER_OPERATOR_ADMIN_USER=root
	export VCLUSTER_OPERATOR_ADMIN_PASSWORD=admin
	export VCLUSTER_OPERATOR_SECRET_KEY="secret key"
	go run main.go server

.PHONY: all docker deploy server
