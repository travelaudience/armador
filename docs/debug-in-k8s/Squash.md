# Squash

https://squash.solo.io

## Pre-Installation

NOTE: rbac permissions per developer are needed. So something like: _(contact someone in devops team for help with this)_
```
kubectl create clusterrolebinding <USER_NAME>-cluster-admin-binding --clusterrole=cluster-admin --user=<USER_NAME>@travelaudience.com
```

## cli App
Install from instructions in release page: https://github.com/solo-io/squash/releases
 and add it to your $PATH
- useful: `squashctl completion`

### VS Code
1. Extenstions -> Install
2. Preferences -> Settings -> (search for squash config) -> (add path to `squashctl`)

### Intellij
_maybe add extension, not sure if it actually does anything though..._
https://plugins.jetbrains.com/plugin/10397-squash-debugger-extension


# Usage

### CLI
1. `squashctl`
2. Choose a debugger to use, a namespace, a pod, and a container to debug
3. When these values have been selected, Squash opens a debug session in you terminal

### VS Code
1. Extenstions -> Install
2. Preferences -> Settings -> (search for squash config) -> (add path to `squashctl`)
3. `CTRL + SHIFT + P` to use the `squash debug` commands
4. select pod
5. select `dlv`


### Intellij
1. Setup a `Go Remote` debug config
	- see screenshot from 8th May
	- set it so that it displays every time, because the port will change
2. run `squashctl` in terminal
3. copy the port when it starts up: `Handling connection for ____`
4. start debugger in IDE and add this port


# Cleanup:
1. squashctl utils list-attachments
2. .... hmm, looks like it's all or nothing - probably not good if someone else is debugging


# Delve

this tool needs it's own training, and I'm not capable of that

Mailing list: https://groups.google.com/forum/#!forum/delve-dev
Best issue ever: https://github.com/go-delve/delve/issues/1537

-------------

Need to build go binary without optimization: `-gcflags "-N -l"`
use ksync to copy it to the pod: (see ksync readme)
