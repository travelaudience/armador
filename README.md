# Armador

A tool for creating ephemeral test environments in Kubernetes.

The best analogy for what Armador provides would be to compare it to Docker Compose, but for Kubernetes.

[Helm](https://helm.sh/) is the building block that `Armador` is based upon. The idea is that a small configuration file is provided for core dependencies for the environment, and then in each chart that needs another chart to be built, that gets referenced inside a `armador.yaml` file in the chart.

### Why not use an "umbrella" helm chart?

The limitation with umbrella charts is that the individual services don't have their own lifecycle. The other problem is that the developer of one service doesn't know all the other dependencies another service has. So they are not able to build their own umbrella charts on demand.

## Gettings started

**Building the code**

```
go get ./...
go build -o $GOPATH/bin/armador
```

```
armador --help
```

### Adding a global configuration file

You will need a configuration file to help the tool know what to install. The default should be located in your $HOME path: `~/.armador.yaml`, however you can use `--config` to point to any other file of your choosing.

An example of the content for your file:

```
prereqCharts:
  - pubsub:
      repo: stable
  - cert-manager:
      packaged: false
      repo: git@github.com:jetstack/cert-manager.git
      pathToChart: deploy/charts/cert-manager
      overrideValueFiles:
        - values.yaml

additionalValues:
  - repo: git@github.com:YOUR_ORG/REPO_WITH_VALUES_FILES.git
    path:
      - release-vals/values.yaml

cluster:
  google:
    name: clusterName
    zone: europe-west3-b
    project: gcp-project
```

### Adding a configuration to your app

In order for Armador to know what to install and what dependencies your application has, you need to add a config file named `armador.yaml` in the helm chart path of your application _(at the same level as `Chart.yaml`)_. One of the basic settings that you could apply would be something like:
```
# version: 1.0
dependencies:
  - chartName:
      repo: stable
overrideValueFiles:
  - values-dev.yaml
```

## Check tools are configured correct

There's some more work that can be done to automate the validation that your system is ready to work with armador, but a good first start is running:
```
armador helm check
```
to be sure there are no errors, and that all the correct plugins are installed.

For example:

* helm version of client should be the same as server
* `helm diff` plugin should be installed _(if it's not `helm plugin install https://github.com/databus23/helm-diff`)_

## Create

The main goal of the tool is to focus on a single application that needs development (eg: `myService1`). If you keep the code for `myService1` in `/code/path/myService1`, and the helm chart for that service is located within that path, then to use Armador, you would run:
```
$ cd /code/path/myService1
$ armador create YOUR_NAME
```
or alternatively you can choose the path of the app to install:
```
armador create YOUR_NAME -p /code/path/myService1
```
#### After `create`

`kubectl get pods --namespace=YOUR_NAME`
will show you the current state of your env.

And now you can use `kubectl` and `helm` in this namespace as you like. Or simply make the changes and run `armador update YOUR_NAME`

## Contributing

Contributions are welcomed! Read the [Contributing Guide](CONTRIBUTING.md) for more information.

## Licensing

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
