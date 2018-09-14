---
title: "Self-Hosting II: Hosting Git"
date: 2018-08-17T23:57:33-07:00
draft: false
---

As a software developer, your online git repository plays a really large role
in your life. It ends up playing part social media, part professional resume,
part file server. It probably hosts a precious `.bashrc` or `.vimrc` file you
can't live without and too many pipe dreams to count. 

When I started to self-host some of the online services I depend on, this was a
natural first target. It's a service that is close to my heart, but it's also
one that would help me in the whole self-hosting process; I could use it as a
place to hold my Kubernetes configurations and other self-hosting related code. 

Thankfully, the quality of free and open source web based git repositories is really high. Here is a short list of the choices available:

- [Gitlab](https://about.gitlab.com/): *Enterprised solution with a open-source core you can host yourself*
- [Gogs](https://gogs.io): *Awesome free and open source git server. Ships as a simple standalone binary*
- [Gitea](https://gitea.io/en-us/) *Fork of Gogs. Seems to be under more active development*

As I only planned to host my own account / code, Gogs or Gitea were the clear
choices. Both are light on resources and super simple to deploy. Gitlab, on the
other hand, seemed all too happy to eat a mountain of RAM for a modestly larger
feature set. I also like that they are community driven so there is no divide
between 'community' and 'entreprise' features. Of the two, I chose Gitea more
or less arbitrarily (\*cough\* dark theme \*cough\*).

## Implementation details

As mentioned in a [previous article](/articles/self-hosting-part-1/), I'm
hosting all my services on a single node Kubernetes 'cluster'. Here is a quick
walk through of deploying Gitea using Kubernetes.

A web based git server hosts two services: the web service and an ssh service.
Once we have Gitea running, we'll need to make these two services accessable
from the outside world. Gitea curates an official Docker image `gitea/gitea`
that you can use to build a deployment. By default the web service on this
image is running on port 3000 and the ssh service is running on port 22.

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gitea
  labels:
    app: gitea
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gitea
  template:
    metadata:
      labels:
        app: gitea
    spec:
      containers:
      - name: gitea-server
        image: gitea/gitea
        ports:
        - containerPort: 22    # SSH service
        - containerPort: 3000  # HTTP service
        volumeMounts:
        - mountPath: /data
          name: gitea-volume
      volumes:
      - name: gitea-volume
        hostPath:
          path: /data/gitea
          type: DirectoryOrCreate
```

Once you have this deployment running on the cluster, its time to configure the
system before exposing it to the outside world. This can be done by
port-forwarding the web service to you local computer.

```shell
$ kubectl port-forward gitea-xxxxxxxxx-xxxxx 3000:3000
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
.
.
```

Open a browser to http://localhost:3000 and you'll be able to create and admin
user and configure some other details like the host name. 

<center>
<img width="80%" src=/img/gitea-config.png></img>
</center>

On set-up you can also disable self-registration. This is handy if you don't
want other people to make accounts on your git server without your permission.
You'll also need to have a little forethought here about your ssh access. I
opted to set up ssh access using a NodePort on port 30022 so I needed to set
SSH Server Port to 30022 on initial configuration. Don't worry about messing
this up. You can also change your configuration later by 'exec'ing into the pod
and manipulating the config with `kubectl exec -it gitea-........ /bin/bash`.
The config file is at `/data/gitea/conf/app.ini` by default. 

Next, we'll set up http access from the outside world. To do this we first
create an service in the cluster that exposes port 3000 of the deployment and
then we point an ingress resource to that service.

```yaml
# http-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: gitea-http
  labels:
    app: gitea
spec:
  type: ClusterIP
  selector:
    app: gitea
  ports:
  - name: http
    port: 3000
    protocol: TCP
```

```yaml
# ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: gitea
spec:
  tls:
  - hosts:
    - code.nfsmith.ca
    secretName: tls-certs
  rules:
  - host: code.nfsmith.ca
    http:
      paths:
      - path: /
        backend:
          serviceName: gitea-http
          servicePort: 3000
```

Note the tls section of the ingress resource. If you populate a tls secret on
your Kubernetes cluster with certificates from Lets Encrypt, you can secure web
access to your git server using TLS! For more details on setting up Ingress,
check out the [last article](/articles/self-hosting-part-1) on setting up
Kubernetes for self hosting. 

Finally, set up access to your SSH service. Like I mentioned before, I opted to
do this with a NodePort service, put you could also do this by stipulating a
host port in the deployment.

```yaml
# ssh-service.yaml
apiVersion: v1
kind: Service
metadata:  
  name: gitea-ssh
  labels:
    app: gitea
spec:
  type: NodePort
  selector:    
    app: gitea
  ports:  
  - name: ssh
    port: 22
    targetPort: 22
    nodePort: 30022
    protocol: TCP
```

Voila, with the NodePort service deployed you should be ready to start using
your brand new git server. As with other deployments its worth checking you can
reboot your server and that the git server comes back up on reboot as you
expect. I've found that Gitea is extremely reliable against this kind of abuse
(I can't say the same for my email unfortunately..)

## Final Thoughts

I've been happily using Gitea as my main git repository for about 2 months now
with no problems. I've also been able to connect a CI/CD service to it using
[Drone](https://drone.io) which I'll likely write about in a future article.
The one role that a self-hosted git server is bad at playing is the social
component. Unfortunately there don't seem to be any git servers that support
any type of federation so it is quite difficult to collaborate with other
servers. There has been some discussion around this for Gitea, but it is
competing against some higher priority work for the time being. Hopefully we'll
see some improvement in this type of thing in the next couple years.

