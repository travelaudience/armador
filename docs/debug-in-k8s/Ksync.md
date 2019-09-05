# Ksync

* https://vapor-ware.github.io/ksync/
* Install: https://github.com/vapor-ware/ksync#installation

> Requires `ksync init` [once] in the cluster to setup the DaemonSet that runs across all namespaces.

On your local dev machine, run
```
ksync watch
```
in a new terminal window to start the local client.

* Use `ksync get` to check the status of what you're currently watching
* Use `ksync delete [NAME]` to remove any current connections that are setup

## Connecting a go app

* Get the app name using this command to see the labels. The app name will be used in the following commands
```
kubectl get pods -n NAMESPACE --show-labels
```

* For a go app, it makes sense to start by syncing the pods directory with an empty dir on your local machine.
This way the currently running binary will be available locally,
and you can configure the `make k8s-debug` step to copy the changed binary to this path.
```
mkdir -p ~/.armador/ksync/APP_NAME
rm -r ~/.armador/ksync/APP_NAME/*
```

* Create the connection
```
ksync create --selector=app=APP_NAME --name=APP_NAME -n NAMESPACE ~/.armador/ksync/APP_NAME /usr/local/bin/
```

-----

You should now have the sync connection created. So to update your app, you need to:
* build a new binary of your changes
* copy the binary to the ksync folder created above
* and wait a few seconds for the file to transfer onto the pod

or just add something like this to your makefile:
```
k8s-debug: build
	cp bin/___ ~/.armador/ksync/APP_NAME/
```

-----

NOTES FOR HOW THINGS COULD WORK BUT DON'T

* https://blog.jetbrains.com/go/2018/04/30/debugging-containerized-go-applications/
* https://blog.openshift.com/debugging-java-applications-on-openshift-kubernetes/
* https://radu-matei.com/blog/state-of-debugging-microservices-on-k8s/
