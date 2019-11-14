# octant and antrea-octant-plugin installation

## Overview

There are two ways to deploy octant and antrea-octant-plugin.

* Deploy octant and antrea-octant-plugin as a process (out of Kubernetes cluster).

* Deploy octant and antrea-octant-plugin as a Pod (inside Kubernetes cluster).

### 1. Deploy octant and antrea-octant-plugin as a process

See details in [octant installation](https://github.com/vmware-tanzu/octant/blob/master/README.md#installation)

You can follow the steps listed below to install octant and antrea-octant-plugin on linux.

1.Get and install octant v0.8.0.
Depend on your linux operation system, to install octant v0.8.0, you can use either

```
wget https://github.com/vmware-tanzu/octant/releases/download/v0.8.0/octant_0.8.0_Linux-64bit.deb

dpkg -i octant_0.8.0_Linux-64bit.deb
```
or
```
wget https://github.com/vmware-tanzu/octant/releases/download/v0.8.0/octant_0.8.0_Linux-64bit.rpm

rpm -i octant_0.8.0_Linux-64bit.rpm
```

2.Export your kubeconfig path (file location depends on your setup) to environment variable $KUBECONFIG.

```
export KUBECONFIG=/etc/kubernetes/admin.conf
```

3.May need to compile antrea to get antea-octant-plugin if you did not "make bin".

```
make bin
```

4.Move antrea-octant-plugin to OCTANT_PLUGIN_PATH.

```
# If you did not change OCTANT_PLUGIN_PATH, the default folder should be $HOME/.config/octant/plugins.

mv antrea/bin/antrea-octant-plugin $HOME/.config/octant/plugins/
```

5.Start octant as a backstage process with UI related environment variables.

```
#  Change port 8900 according to your environment.

OCTANT_LISTENER_ADDR=0.0.0.0:8900 OCTANT_ACCEPTED_HOSTS=$HOSTNAME OCTANT_DISABLE_OPEN_BROWSER=true nohup octant &
```

Now, you are supposed to see octant is running together with antrea-octant-plugin via url http://$HOSTNAME:8900.


### 2. Deploy octant and antrea-octant-plugin as a Pod

For deploying octant in Pod, please refer to the example [Running octant in cluser](https://github.com/vmware-tanzu/octant/tree/master/examples/in-cluster)

Following the example above, you need to execute the commands below in order to running octant and antrea-octant-plugin in Pod.

1. Build antrea-octant docker image.

```
make antrea-octant-ubuntu
```

2.Create a secret that contains your kubeconfig.

```
# Change --from-file according to kubeconfig location in  your set up.

kubectl create secret generic octant-kubeconfig --from-file=/etc/kubernetes/admin.conf -n kube-system
```

3.May need to update build/yamls/antrea-octant.yml according to your kubeconfig file name.

4.May need to label certain node if you want to run octant Pod on a fixed node so that you can access UI via a fixed URL.

```
kubectl label nodes $HOSTNAME role=antrea-octant.
```

5.Apply the deployment.

```
kubectl apply -f build/yamls/antrea-octant.yml
```

Now, you are supposed to see octant is running together with antrea-octant-plugin via url http://$HOSTNAME:8900.
