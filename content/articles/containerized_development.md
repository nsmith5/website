---
title: "Containerized Development Tools"
date: 2019-06-29T23:59:14-07:00
draft: false
---

If you're like me your process for vetting third-party packages for your code
goes like this:

- Search for "how to do {problem of the moment} in {language of the day}"
- 10 points for a Github page that is in an organization instead of personal
- 10 points per cute colourful badge on README
- 30 points for an actual website
- 1 point per Github star
- basic multiplier for my desperation level

If I find a winner, `go get` that sucker (or what ever the package manager of
the day is) and try it out. If it solves the problem then I _should_ go vet the
code I'm running. What do I actually do? Recognizing its a problem I've mostly
just developed a prejudice against langauges that don't have a good, broad
standard library and try to keep my dependency tree as small as possible.

There is are a mountain of problems with this. Even if I _did_ vet my
dependencies after the fact, the package has a chance to run a bunch of code at
install time on my machine. Now, obviously, the best practice would be to read
all of the code before installing it, but on the trade off of security and
convenience, this is rough on convenience.

Security considerations aside, its also annoying that development dependencies
just keep piling up in your system. Every package manager feels like it takes
free reign over your home directory. In a system backup you end up copying
mountains of source code from these dependencies unless you filter them out.
Some of them aren't even gracious enough to restrict their actions to hidden
files (ehem.. I'm looking at you Go).

## Containerize!

One way to get around this is to containerize your dev tools. The idea is to
mount your project directory into a container.

```shell
$ podman run -it --rm -v $PWD:/workspace:z golang /bin/sh
(in-container) $ cd /workspace
(in-container) $ go build  # Do your golang stuff
(in-container) $ CTRL-D
$
```

The `:z` relabels files so that they appear to be owned by root in the container
and my user outside. Also you can use `docker` in place of `podman` if thats your
thing.

This is good, but not great. Any dependencies you install get nuked every time you
leave the container. Lets fix that by attaching a named volume.

```shell
$ podman run -it --rm -v golang:/go -v $PWD:/workspace:z golang /bin/sh
```

So what have we won here? Dependencies can only read the files you're working
on (not that big a deal because there are no secrets in your repository right?)
and the other dependences you've installed. Super cool right? Unfortunately,
the process is a little rough from a usability perspective. It would be better if
we could just run `go build` and this stuff happened behind the scenes.

Fortunately, thats not too tricky to achieve. Lets change the image a bit:

```Docker
FROM golang:alpine
WORKDIR /workspace
ENTRYPOINT ["/usr/local/go/bin/go"]
```

build it:

```shell
$ podman build -t go /path/to/directory/of/Dockerfile
```

and make a bash alias:

```shell
$ alias go="podman run --rm go:/go -v $PWD:/workspace:z go"
$ go version                         # Victory!
go version go1.12.6 linux/amd64
```

Excellent, looks like the normal tools up front, but isolated dependencies in
the back!

I've been pretty happy with this workflow, so I've started to compile some
container recipes and aliases for this purpose in a repository. If you're
interested in these tools take a look at
[code.nfsmith.ca/nsmith/containers](https://code.nfsmith.ca/nsmith/containers).
