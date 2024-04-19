OTEL_OPERATOR_VERSION=v0.97.0
VERSION=v1

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

.PHONY: manifests-${VERSION}
manifests-${VERSION}: lint-${VERSION}
	-rm -rf $@
	helm template local-release ./chart-${VERSION} \
		--set global.clusterName=sushi \
		--namespace=collection \
		--debug \
		--output-dir $@

validate-0: manifests-0 opentelemetry.io_opentelemetrycollectors.yaml
	 kubectl-validate ./manifests-0 --local-crds .

.PHONY: lint-${VERSION}
lint-${VERSION}:
	helm lint ./chart-${VERSION} --set global.clusterName=ci

.PHONY: install-${VERSION}
install-${VERSION}:
	helm upgrade --install --set global.clusterName=local test --create-namespace -n collection ./chart-${VERSION}

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

unit-tests-${VERSION}:
	$(MAKE) -C ./chart-${VERSION} unit-tests

unit-tests-config:
	$(MAKE) -C ./chart/test-config unit-tests

