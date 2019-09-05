# More info about using Armador

In this directory, you'll find:

* [example](example): a close to "real-world" use case of implementing Armador
* [debug-in-k8s](debug-in-k8s): Documentation on tips to use tools like Telepresence/Kysnc/Squash for debugging an app that's running in an environment created by Armador

## Yaml files defined

The "global config" (aka: `~/.armador.yaml`)

```yaml
prereqCharts:                   # The charts that will be installed in any env you create. This is useful for things like databases, proxies, etc. It accepts a list of charts to install. (optional)
  - postgres:                   # This is the name of the Helm chart. It needs to match exactly.
      repo: stable              # [required] The Helm repo to pull the chart from. If packaged=false, than this would be the git repo (see next example).
  - cert-manager:
      repo: git@github.com:jetstack/cert-manager.git
      packaged: false           # [optional] Defines if Armador should treat the chart as a tarball, or a directory of yaml. (default true)
      pathToChart: deploy/charts/cert-manager   # [optional] If packaged=false, where in the repo is the chart located.
      overrideValueFiles:       # [optional] A list of files that will override values.yaml. Relative to the path of the chart.
        - values.yaml

additionalValues:               # A list that consists of a universal set of values that override all other retrieved settings. Useful for setting things in prereqCharts that come from public charts that can't be modified. (optional)
  - repo: git@github.com:YOUR_ORG/REPO_WITH_VALUES_FILES.git  # [required] Where to obtain the values from
    path:                       # A list of files to apply. Relative to within the repo.
      - release-vals/values.yaml

cluster:                        # [required] The configuration for what k8s cluster to create envs in. (currently only supporting GKE & minikube. If both are set, GKE settings take precendence.)
  google:                       # If using GKE, use this key (the following values are then required)
    name: clusterName
    zone: europe-west3-b
    project: gcp-project
  minikube:                     # If connecting to minikube, use this key (the following values are then required)
    contextName: minikube
```

The individual `armador.yaml` (in each chart).

```yaml
dependencies:                   # [optional] A list of Helm charts that this application requires to be "functional".
  - chartName:                  # These follow the same structure as `prereqCharts` in the global config (see above).
      repo: stable
overrideValueFiles:             # [optional] A list of files that will override values.yaml for THIS app. Relative to the path of the chart.
  - values-dev.yaml
```
