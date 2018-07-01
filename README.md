# My Website 
*www.nfsmith.ca*

Source code for my website. Contains a static hugo site and a simple static file
server written in go to serve it. To build the server run,

```shell
$ CGO_ENABLED=0 go build -o server -ldflags '-s' src/main.go
```

To build the website run,

```shell
$ hugo
```

To packge the whole thing up in a container, build the server and website and
run,

```shell
$ docker build -t website .
$ docker run -p 3000:3000 website
```
