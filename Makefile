OTEL_OPERATOR_VERSION=v0.97.0

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

install-kubeconform:
	go install github.com/yannh/kubeconform/cmd/kubeconform@latest

opentelemetry.io_opentelemetrycollectors.yaml:
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-operator/${OTEL_OPERATOR_VERSION}/bundle/manifests/opentelemetry.io_opentelemetrycollectors.yaml

opentelemetry.io_opentelemetrycollectors.json: opentelemetry.io_opentelemetrycollectors.yaml
	cat $< | yq -o=json '.spec.versions[0].schema.openAPIV3Schema' > $@


kubeconform: manifests opentelemetry.io_opentelemetrycollectors.json
	kubeconform \
		-strict \
		-summary \
		-schema-location opentelemetry.io_opentelemetrycollectors.json \
		manifests

unit-tests-0:
	$(MAKE) -C ./chart-0/test-templates unit-tests

unit-tests-config:
	$(MAKE) -C ./chart/test-config unit-tests

