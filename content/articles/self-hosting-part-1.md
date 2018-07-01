---
title: "Self Hosting Part 1: Setting up Kubernetes"
date: 2018-06-30T14:13:28-07:00
draft: false
---

This article is the first in a series about my experience moving away from 
free, hosted webservices (email, git, social, file storage etc) and towards 
self-hosting. To host my services I decided to use a stand alone Kubernetes
'cluster'. 

I chose Kubernetes for the same reason it was orginally made: I want to pack as 
many services as possible on one server. There is some overhead from the 
Kubernetes daemons (etcd, the API server, kubelet etc), but once in place its a 
great system for deploying heaps of services and effeciently routing traffic to 
them.

## Choosing a Server and Hosting

I use [Fedora](https://getfedora.org/en/workstation/) as my daily operating
system to choosing it for my server operating system was an obvious choice. 
Fedora has an flavour specifically for deploying containers called 
[Fedora Atomic](https://getfedora.org/en/atomic/) but, I've found the 
documentation a little weak and its future isn't that clear to me following the
[CoreOS acquisition](https://www.redhat.com/en/about/press-releases/red-hat-acquire-coreos-expanding-its-kubernetes-and-containers-leadership), 
so I opted for the traditional server spin. For hosting, I opted for 
[Digital Ocean](https://www.digitalocean.com/) and choose the 2 CPU / 4Gb RAM
option.

Setting up automatic updates is super easy on Fedora so I did that first,

```shell
server $ dnf update && systemctl reboot // Update server
server $ dnf install dnf-automatic
server $ dnf enable --now dnf-automatic-install.timer
```

## Installing Kubernetes

[Kubeadm](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/) 
is an excellent tool for spinning up a Kubernetes cluster. Installing on Fedora
is fairly straight forward. You can follow the tutorial [here](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/),
but I'll show the detailed steps for clarity. There are a few Fedora specific
details worth noting. 

To begin, install Kubeadm from the repositories and enable Docker.

```shell
server $ dnf install kubernetes-kubeadm docker
server $ systemctl enable --now docker
```

You need to set selinux to permissive as well.

```shell
server $ setenforce 0  // Temporary
server $ vi /etc/sysconfig/selinux
# This file controls the state of SELinux on the system.
# SELINUX= can take one of these three values:
#     enforcing - SELinux security policy is enforced.
#     permissive - SELinux prints warnings instead of enforcing.
#     disabled - No SELinux policy is loaded.
SELINUX=permissive      # <----- Set to 'permissive'
# SELINUXTYPE= can take one of these three values:
#     targeted - Targeted processes are protected,
#     minimum - Modification of targeted policy. Only selected processes are protected.
#     mls - Multi Level Security protection.
SELINUXTYPE=targeted
server $
```

Once done, we're ready to deploy Kubernetes! I decided to use Flannel for my 
networking and pointed out where that choice effected certain commands. If you
choose a different network lookup the relevant values in the kubeadm docs.

```shell
server $ kubeadm init --pod-network-cidr=10.244.0.0/16 # <-- Important for setting up flannel network later
server $ sysctl net.bridge.bridge-nf-call-iptables=1 # <-- Also import for flannel
server $ export KUBECONFIG=/etc/kubernetes/admin.conf
server $ kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/v0.10.0/Documentation/kube-flannel.yml # Deploy flannel!
```

Finally, we need to remove a taint on the master node to allow pods to run 
there.

```shell
server $ kubectl taint nodes --all node-role.kubernetes.io/master-
```

To connect to your cluster remotely copy `/etc/kubernetes/admin.conf` from the
server and paste it into `/.kube/config` on your local machine.

```shell
laptop $ kubectl get nodes
NAME        STATUS    ROLES     AGE       VERSION
server      Ready     master    3d        v1.10.1
```

## Setting up Ingress

We have one last thing to set up before our cluster is ready to use: ingress.
An ingress controller routes traffic to services based on the domain name and
path requested. This is great because it means we can publish many services on
the same server and distinguish them by domain.

If you have a domain name available, you can set this up by making a wild card
record for you domain and pointing it towards your server. For example, my
domain is nfsmith.ca so I set up a record for *.nfsmith.ca and directed that
towards my droplet. This means I can host any subdomain of nfsmith.ca at the 
droplet (www.nfsmith.ca, code.nfsmith.ca etc etc).

The ingress controller itself is very easy to install. You can choose from a 
variety of ingress controllers. I'm familiar with HAProxy so I choose one based
on that. Deploying is as simple as,

```shell
laptop $ kubectl apply -f kubectl create -f https://raw.githubusercontent.com/jcmoraisjr/haproxy-ingress/master/docs/haproxy-ingress.yaml
laptop $ kubectl label node server role=ingress-controller
```

Checkout out the [docs](https://github.com/jcmoraisjr/haproxy-ingress) for more 
information on configuration.

## Finishing up

At this point its best to reboot the system and make sure everything comes back
up afterwords. If all goes well, you're ready to start hosting services. Next
time we'll take a look at hosting our own git repostory and securing services
with LetsEncrypt certificates.