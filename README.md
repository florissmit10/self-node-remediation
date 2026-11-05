# Self Node Remediation - Automatic Kubernetes Node Remediation 
<p align="center">
<img width="200" src="config/assets/snr_icon_blue.png">
</p>

Existing baremetal remediation strategies utilize BMC credentials to power-cycle and/or reprovision the host.
However there are also environments that either do not include BMCs, or there are policies
in place that prevent them from being utilized.  Such environments would also benefit from
the ability to safely recover affected workloads and restore cluster capacity (where possible).
This self node remediation controller is using an alternate mechanism for a node in a cluster to detect its health
status and take actions to remediate itself in case of a failure.  While not all remediation events can
result in the node returning to a healthy state, the proposal does allow surviving parts of the cluster
to assume the node has reached a safe state so that it’s workloads can be automatically recovered.
This work can also be useful for clusters with BMC credentials.


## More Info
https://www.medik8s.io/

## Project State
The operator is available in [operator hub](https://operatorhub.io/operator/self-node-remediation)

Self Node Remediation works best together with [Node Health Check Operator](https://operatorhub.io/operator/node-healthcheck-operator)

## Installation without OLM

On Kubernetes without OLM, SNR can be installed with the `install.yaml` manifest attached to the
[releases](https://github.com/medik8s/self-node-remediation/releases). Older releases don't have
this asset; install those with OLM. It needs
[cert-manager](https://cert-manager.io/docs/installation/) for the webhook serving certificate.

```bash
kubectl apply --server-side -f https://github.com/medik8s/self-node-remediation/releases/download/v<version>/install.yaml
```

SNR is installed into the `self-node-remediation` namespace. The operator runs its agent as a
privileged DaemonSet in that namespace, so the manifest labels the namespace with the `privileged`
[Pod Security Standard](https://kubernetes.io/docs/concepts/security/pod-security-admission/).

To build the manifest from source, run `make build-installer VERSION=<version>`, with the version
without the leading `v`, which writes `dist/install.yaml`. Use `IMG` to set a different image; it is
used for both the operator and the agent.

### Using Helm

The Helm chart in `dist/chart` installs the same resources and also needs cert-manager.
Helm does not label the namespace, so create it with the `privileged` Pod Security level first.
Run this from a checkout of the repository at a release tag that includes the chart:

```bash
kubectl create namespace self-node-remediation
kubectl label namespace self-node-remediation \
  pod-security.kubernetes.io/enforce=privileged \
  pod-security.kubernetes.io/audit=privileged \
  pod-security.kubernetes.io/warn=privileged
helm install self-node-remediation dist/chart --namespace self-node-remediation
```

The agent uses the operator image unless `agent.image` is set. See `dist/chart/values.yaml` for the
configurable values. The CRDs are kept when the chart is uninstalled, so existing resources survive
a reinstall.

The chart is generated from `dist/install.yaml` with the kubebuilder Helm plugin.
After changing `config/`, run `make helm-chart` and update `dist/chart/values.yaml` until
`make verify-helm-chart` passes.

## Help
Feel free to join our google group to get more info - https://groups.google.com/g/medik8s
