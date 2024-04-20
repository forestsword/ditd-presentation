OTEL_OPERATOR_VERSION = v0.97.0
VERSION ?= v1

default: test-${VERSION}

.PHONY = create-cluster repos install-local-operator
create-cluster:
	kind create cluster

repos:
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	helm repo update

.PHONY: test-${VERSION}
test-${VERSION}: lint-${VERSION} kubeconform-${VERSION} validate-${VERSION} unit-tests-${VERSION}

install-local-operator:
	helm install \
		opentelemetry-operator open-telemetry/opentelemetry-operator \
		--set admissionWebhooks.certManager.enabled=false \
		--set admissionWebhooks.autoGenerateCert.enabled=true

.PHONY: manifests-${VERSION}
manifests-${VERSION}: lint-${VERSION}
	-rm -rf $@
	helm template local-release ./chart-${VERSION} \
		--set global.clusterName=sushi \
		--namespace=collection \
		--debug \
		--output-dir $@

validate-${VERSION}: manifests-${VERSION} opentelemetry.io_opentelemetrycollectors.yaml
	 kubectl-validate ./manifests-${VERSION} --local-crds .

.PHONY: lint-${VERSION}
lint-${VERSION}:
	helm lint ./chart-${VERSION} --set global.clusterName=ci

.PHONY: install-${VERSION}
install-${VERSION}:
	helm upgrade --install \
		--set global.clusterName=local \
		--create-namespace \
		-n collection \
		test \
		./chart-${VERSION}

install-kubeconform:
	go install github.com/yannh/kubeconform/cmd/kubeconform@latest

opentelemetry.io_opentelemetrycollectors.yaml:
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-operator/${OTEL_OPERATOR_VERSION}/bundle/manifests/opentelemetry.io_opentelemetrycollectors.yaml

kubeconform-${VERSION}: manifests-${VERSION}
	kubeconform -ignore-missing-schemas manifests-${VERSION}

unit-tests-${VERSION}:
	$(MAKE) -C ./chart-${VERSION} unit-tests
