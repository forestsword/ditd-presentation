OTEL_OPERATOR_VERSION = v0.97.0
VERSION ?= v1

default: test

.PHONY: create-cluster
create-cluster:
	-kind create cluster

.PHONY: repos
repos:
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	helm repo update

.PHONY: test
test: lint kubeconform validate unit-tests

.PHONY: install-local-operator
install-local-operator: repos create-cluster
	helm install --upgrade \
		opentelemetry-operator open-telemetry/opentelemetry-operator \
		--set admissionWebhooks.certManager.enabled=false \
		--set admissionWebhooks.autoGenerateCert.enabled=true

.PHONY: manifests
manifests: lint
	-rm -rf $@
	helm template local-release ./chart-${VERSION} \
		--validate \
		--set global.clusterName=sushi \
		--namespace=collection \
		--debug \
		--output-dir $@

.PHONY: validate
validate: manifests opentelemetry.io_opentelemetrycollectors.yaml
	 kubectl-validate ./manifests-${VERSION} --local-crds .

.PHONY: lint
lint:
	helm lint ./chart-${VERSION} --set global.clusterName=ci

# This is NOT what ArgoCD does:
# https://argo-cd.readthedocs.io/en/stable/faq/#after-deploying-my-helm-application-with-argo-cd-i-cannot-see-it-with-helm-ls-and-other-helm-commands
.PHONY: install
install:
	helm upgrade --install \
		--set global.clusterName=local \
		--create-namespace \
		-n collection \
		local \
		./chart-${VERSION}

.PHONY: uninstall
uninstall:
	helm uninstall \
		-n collection \
		local 

install-kubeconform:
	# Requires go installed and configured for your PATH
	go install github.com/yannh/kubeconform/cmd/kubeconform@latest

opentelemetry.io_opentelemetrycollectors.yaml:
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-operator/${OTEL_OPERATOR_VERSION}/bundle/manifests/opentelemetry.io_opentelemetrycollectors.yaml

.PHONY: kubeconform
kubeconform: manifests
	kubeconform -ignore-missing-schemas manifests-${VERSION}

.PHONY: unit-tests
unit-tests:
	$(MAKE) -C ./chart-${VERSION} unit-tests
