# Telepresence

## Connecting a go app

* https://www.telepresence.io/reference/connecting

* only one telepresence can run per machine, and you _can't_ use other VPNs.
* DO NOT use this in `staging` or `production`

Create a new telepresence deployment that will route traffic from your deployment to localhost
```
telepresence --namespace <YOUR_NAMESPACE> --swap-deployment <DEPLOYMENT_NAME> --run-shell
```

* If you use the `--expose` option for telepresence with a given port the pod will forward traffic it receives on that port to your local process.
* If you have more than one container in the pods created by the deployment you can also specify the container name.
* If telepresence crashes you may need to manually restore the Deployment.
* It is possible to expose a different local port than the remote port, eg `--expose 8080:80`.
* Additional cloud resources can be routed via the cluster to localhost if you explicitly specify them using `--also-proxy`
* Volumes configured in your Deployment pod template can be made available to your local process (best performance with read-only volume mounts)
* You can add `--env-json env-vars.json` to get an ouput of which environment variables are made available.
