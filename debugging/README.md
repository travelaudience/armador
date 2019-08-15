# Debugging running app in Kubernetes

There's a few techniques that can be used and in the end it depends on what the goals are. Regardless, you'll want to make use of different tools that exist for this kind of thing.

**NOTE**: Because most of these techniques introduce latency, and other potential problems to end users of the system, it is highly suggested that you
**DO NOT DO THIS IN PRODUCTION**.

## [Telepresence](Telepresence.md)

The way Telepresence works, is that it creates a proxy from Kubernetes to your local env. So lets say you have two services (A & B) running in K8S and you're developing service A. You use Telepresence to swap the service A deployment, and then any traffic that goes from service B to A will go to your local env. This makes it easy to add quick changes and see how they impact the full system. It also allows for setting breakpoints in your IDE.

PROBLEM: there's a bug that makes testing this out not possible: https://github.com/telepresenceio/telepresence/issues/972

## [Ksync](Ksync.md)

The power of Ksync is that it _syncs_ files from your local env to kubernetes pods. Even if your app is a go binary that is running on a pod, Ksync can copy the changes you make to that binary to the pod, and update the running application with the pod needing to be restarted.

## [Squash](Squash.md)

Used to create debugging connections from local env to running pods.

-----------

Further reading for each tool is available in the respective files.
