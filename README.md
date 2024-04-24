# The Daemonset in the Details

This is the companion code for the presentation
_The Daemonset in the Details_
given at the Kubernetes Community Days (KCD)
Bucharest, Romania April 2024.

## Usage

To actually execute the Helm chart you will need a Kubernetes cluster
and `kubectl` and Helm installed.
[Kind](https://kind.sigs.k8s.io) is a good way to do this locally.

### Local Cluster Setup

You will need [`kind`](https://kind.sigs.k8s.io/docs/user/quick-start/).
Once installed run `make create-cluster`:

```
 ditd-presentation  % make create-cluster
kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.29.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

Install the OpenTelemetry Helm Repo `make repos`:

```
ditd-presentation  % make repos
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
```

Deploy the OpenTelemetry Operator `make install-local-operator`:

```
helm upgrade --install \
        opentelemetry-operator open-telemetry/opentelemetry-operator \
        --set admissionWebhooks.certManager.enabled=false \
        --set admissionWebhooks.autoGenerateCert.enabled=true
Release "opentelemetry-operator" does not exist. Installing it now.
NAME: opentelemetry-operator
LAST DEPLOYED: Wed Apr 24 13:54:56 2024
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
opentelemetry-operator has been installed. Check its status by running:
  kubectl --namespace default get pods -l "release=opentelemetry-operator"

Visit https://github.com/open-telemetry/opentelemetry-operator for instructions on how to create & configure OpenTelemetryCollector and Instrumentation custom resources by using the Operator.
```

You should now see the operator running:

```
 ditd-presentation  % kubectl get pod -n default
NAME                                     READY   STATUS    RESTARTS   AGE
opentelemetry-operator-7fb78c8fb-cc2b2   2/2     Running   0          15s
```

## Errata

### Basic Manifest Validations

#### `kubeconform`

See [`kubeconform`](https://github.com/yannh/kubeconform)

#### `kubectl-validate`

See [kubectl-validate](https://github.com/kubernetes-sigs/kubectl-validate)
for offline/client-side validation.
