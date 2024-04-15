

.PHONY = create-cluster repos install-local-operator
create-cluster:
	kind create cluster

repos:
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	helm repo update

install-local-operator:
	helm install \
  opentelemetry-operator open-telemetry/opentelemetry-operator \
	--set admissionWebhooks.certManager.enabled=false \
	--set admissionWebhooks.autoGenerateCert.enabled=true


.PHONY: manifests
manifests: lint
	-rm -rf $@
	helm template local-release ./chart \
		--set global.clusterName=sushi \
		--api-versions opentelemetry.io/v1alpha1/OpenTelemetryCollector \
		--debug \
		--namespace=collection \
		--output-dir manifests

.PHONY: lint
lint:
	helm lint ./chart --set global.clusterName=ci

.PHONY: install
install:
	helm upgrade --install --set global.clusterName=local test --create-namespace -n collection ./chart
