---
title: "Self Hosting Part 1: Setting up Kubernetes"
date: 2018-06-30T14:13:28-07:00
draft: false
---

This article is the first in a series about my experience moving away from 
free, hosted webservices (email, git, social, file storage etc) and towards 
self-hosting. I wrote a bit about my motivation to take this project on in 
[Part 0](/articles/self-hosting-part-0) of the series. You may be interested
in this series if you share those concerns or if you're just looking for a 
fun project to hack on.

The first step to self hosting is *hosting*. I'm a fan of Linux containers
for running services. I've had some chances to use Kubernetes professionally
and really enjoyed the experience. This is motivated by a few things. The first
one is that I'm not *that* good at administering linux systems. Every 
applications has its own config and stores files in different locations. If you
go through the install and uninstall cycle enough times things just get messy.
You can probably chalk this up to my lack of understanding of the Linux file 
system heirarchy and where to go hunting for things. Anyways, with Kubernetes
I just have to worry about running Kubernetes. Application state is bundled up 
in the container for the most part.

That said, running Kubernetes isn't trivial so its worth taking a look at the
installation process in a bit of detail.

## Choosing a Server and Hosting

I use [Fedora](https://getfedora.org/en/workstation/) as my daily operating
system so it made a natural choice as a server operating system. Kubernetes
development is moving quite rapidly now so it can be nice to have a distribution
that keeps up. I've had no trouble upgrading from release to release with the
`dnf system upgrade` tool so the short support life is no obstacle. 

Fedora has an flavour specifically for deploying containers called 
[Fedora Atomic](https://getfedora.org/en/atomic/) but, I've found the 
documentation a little weak and its future isn't that clear to me following the
[CoreOS acquisition](https://www.redhat.com/en/about/press-releases/red-hat-acquire-coreos-expanding-its-kubernetes-and-containers-leadership), 
so I opted for the traditional server spin. 

For hosting, I opted for 
[Digital Ocean](https://www.digitalocean.com/) and choose the 2 CPU / 4Gb RAM
option. The Kubernetes API server alone takes up about 500Mb. 4Gb of RAM may
be excessive for a minimal install, but I think 2 Gb is probably as low as 
you comfortably go. Some services like email can be very RAM heavy so but others
consume almost nothing so it worth thinking about what you want to deploy before
you get started.

With a new droplet its always nice to deal with updates now and in the future
using `dnf-automatic`:

```shell
laptop $ ssh root@server.example.com
server $ dnf update && systemctl reboot // Update server
server $ dnf install dnf-automatic      // Automatic updates
server $ systemctl enable --now dnf-automatic-install.timer
```

## Installing Kubernetes

[Kubeadm](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/) 
is an excellent tool for spinning up a Kubernetes cluster. They even have a good
[tutorial](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)
for setting up a stand-alone cluster. Unfortunately there are a few wrinkles to
following the tutorial on a Fedora server. If you're the impatient and trusting
type you can do the following to get started:

```shell
server $ dnf install kubernetes-kubeadm docker tc
server $ systemctl enable --now kubelet
server $ curl https://code.nfsmith.ca/nsmith/kubernetes-resources/raw/branch/master/deploy.sh > deploy.sh
server $ # Edit deploy.sh : replace <node name> with name of host
server $ bash deploy.sh  # Install everything
server $ # Now go to /etc/sysconfig/selinux and set SELINUX=permissive
```

Once finished, you'll have a new Kubernetes cluster with Flannel network and 
HAProxy ingress controller. To connect remotely, make sure your firewall is open 
to part 6443, and copy the kubernetes config from your server 
(`/etc/kubernetes/admin.conf`) to you local machine (`~/.kube/config`).

```shell
laptop $ kubectl get nodes
NAME        STATUS    ROLES     AGE       VERSION
server      Ready     master    3d        v1.10.1
```

If you're the less trusting type (and you should be!), check out the `deploy.sh`
script at [code.nfsmith.ca](https://code.nfsmith.ca/nsmith/kubernetes-resources/src/branch/master/deploy.sh)


## Ingress

An ingress controller routes traffic to services based on the domain name and
path requested. For example, if traffic for www.example.com and code.example.com
are directed to the same node, the ingress controller route traffic for each 
subdomain to a seperate Kubernetes service. This is great because it means we 
can publish many services on the same server and distinguish them by domain. A
typical ingress resource might look something like this:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: website
spec:
  tls:
  - hosts:
    - www.nfsmith.ca       # Host name for TLS
    secretName: tls-certs  # Certificates for TLS termination
  rules:
  - host: www.nfsmith.ca   # Host name to match
    http:
      paths:
      - path: /            # Path to match
        backend:
          serviceName: website  # Service to accept traffic
          servicePort: 3000     # Service port
```

If you have a domain name available, you can set this up by making a wild card
record for you domain and pointing it towards your server. For example, my
domain is nfsmith.ca so I set up a record for *.nfsmith.ca and directed that
towards my droplet. This means I can host any subdomain of nfsmith.ca at the 
droplet (www.nfsmith.ca, code.nfsmith.ca etc).

Ingress controllers are basically reverse proxies managed by Kubernetes so 
different implementations are typically backed by different proxy software 
(Nginx, HAProxy etc). If you'd like to know more about the HAProxy ingress
controller installed by the script checkout the 
[repository](https://github.com/jcmoraisjr/haproxy-ingress) here.

## Finishing up

At this point its best to reboot the system and make sure everything comes back
up afterwords. If all goes well, you're ready to start hosting services. Next
time we'll take a look at hosting our own git repostory and securing services
with LetsEncrypt certificates.
