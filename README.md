# The Daemonset in the Details

This is the companion code for the presentation
_The Daemonset in the Details_
given at the Kubernetes Community Days (KCD)
Bucharest, Romania April 2024.

## Usage

To actually execute the Helm chart you will need a Kubernetes cluster
and Helm installed.
[Kind](https://kind.sigs.k8s.io) is a good way to do this locally.

## Overall Questions?

- What are the guarantees of our tools
  or the technologies we use?
- Can we verify our logic?
  Can we verify what goes in and what comes out is what we expect;
  i.e. can we introspect?

## The Story

### Week 1

API server rejects what Helm says should be valid.

Basic validation: `kubeconform`

### Week 2

API server accepts but the collector crash-loops

Basic validation: `otelcol validate` ... TODO

### Week 3

Do the business - how to avoid constant manual verification

Basic validation: Manual verification after deployment

## Conclusions

- Learn the language of the tools you use
  If the technologies were using Python this code would be in Python
- The otel project is great.
  I don't have to copy code.
  Kudos to helm as well.
- How can otel make this simpler?
  - components
  - is golang the right way?
    Lessons from prom/amtool?
- Give them tests let them eat bread
- Safer stress free deployments
- Why not create your own distribution and component?

## Errata

### `kubeconform`

You need first `openapi2jsonschema`:

```
pip install openapi2jsonschema
```

Convert but shit's broken... TODO

then run `kubeconform`... TODO

### `kubectl-validate`

See [kubectl-validate](https://github.com/kubernetes-sigs/kubectl-validate)
for offline/client-side validation.
