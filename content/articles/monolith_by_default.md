---
title: "Monolith by Default, Microservice with Configuration"
date: 2020-08-17T00:00:00-07:00
draft: false
---

If you look around the web services industry at the moment, you'll see a sort
of thrashing between microservice and monolithic architectures. Wins in the
microservices architecture have been smaller, single purpose services that
are easier to test and deploy. On the other hand, the aggregate service can
become much more difficult to operate and test. Complex observability
tooling, such as distributed tracing, can be needed to understanding system
failures. While you might have blasted out your microservice product in no
time, you might be starting to notice your SRE team is screaming for
resources and has starting becoming a sort of glue team that needs to
understand how every microservice works to debug nasty distributed failures.

Once you start down the microservice pathway, its natural to see organic growth
happen in which more and more services get added on. Caching systems, queues,
bespoke microservices can all start to be added to a growing soup of
interdependency until you get to the point that its just impossible to run your
product on a laptop. Once you're there, the nasty drop off into slow
development iteration is inevitable.

How can we get around this? I think one approach is to insist that your
application is monolithic by default with the option to be broken into
microservices by configuration. I think this is similar to [Peter Bourgon's
insistence][1] that `go test` should always succeed. Just as you should opt
in to stateful testing, you should opt into microservices.

[1]: https://twitter.com/peterbourgon/status/989571449856798720

You should always be able to spin up some version of your product
with one command:

```shell
$ my-app
Listing on default port 8080...
```

Want to split out the authentication duties? Just add some flags

```shell
$ my-app auth --port 8081 &
Auth service listening on port 8081...
$ my-app --auth-service localhost:8081
Listening on default port 8080...
```

## 'Microlith' Example

Lets walk through an example of what that might look like in Go. Our product
is going to be a bossy little service called "microlith" that tells you what
you can and can't do. 

> Note: You can find the source for this example [here][2].

We'll start out with a monolith:

```golang
package main

import (
	"fmt"
	"net/http"
	"strings"
)

var okThings = map[string]struct{}{
	"walk": struct{}{},
	"talk": struct{}{},

	// Nah, we're not that cool
	// "walk the walk": struct{}{}
}

type API int

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	thing := strings.ReplaceAll(r.URL.Path, "/", "")
	_, ok := okThings[thing]
	if ok {
		fmt.Fprintf(w, "Oh yeahh you can really %s\n", thing)
	} else {
		fmt.Fprintf(w, "Uh oh, looks like you can never %s\n", thing)
	}
}

func main() {
	http.ListenAndServe(":8080", new(API))
}
```

It might be nice to break off the permissions logic into something separate
right? We could have a different team work on its code. Maybe its performance
characteristics and scaling properties are different than the rest of the
application and we might be able to optimize its deployment if it ran on its
own. For _whatever_ reason we'd like to break this bit of responsibility out
into something independent.

## Creating Interfaces

Lets start by creating a separate `permissions` package and a public
interface that we can consume.

```golang
# permissions/interface.go
package permissions

type Service struct {
    CanI(string) bool
}
```

And we create a simple implementation that will work in our monolith:

```golang
# permissions/simple.go
package permissions

var okThings = map[string]struct{}{
	"walk": struct{}{},
	"talk": struct{}{},

	// Nah, we're not that cool
	// "walk the walk": struct{}{}
}

type simple int

func (s *simple) CanI(thing string) bool, error {
    _, ok := okThings[thing]
    return o, nil
}

func NewSimple() Service {
    return new(simple)
}
```

We can now refactor our API as follows:

```golang
package main

import (
	"fmt"
	"net/http"
	"strings"

        // Path to _your_ permissions package here
	"github.com/nsmith5/microlith/permissions"
)

type API struct {
	ps permissions.Service
}

func (a API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	thing := strings.ReplaceAll(r.URL.Path, "/", "")
	err, ok := a.ps.CanI(thing)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } 
    if ok {
		fmt.Fprintf(w, "Wow, you can really %s\n", thing)
	} else {
		fmt.Fprintf(w, "Uh oh, looks like you can never %s\n", thing)
	}
}

func main() {
	http.ListenAndServe(":8080", API{permissions.NewSimple()})
}
```

Ok! This is great. This refactor has created a border in our product that two
different teams can work on either side of. This parcelling up of turf
between development teams is most of the microservices battle. You should
_probably_ stop here.

## Microservices behind the interface

So your product has grown over time and you've noticed that the permissions
calls are basically CPU bound and the API is bound by network resources or
something. Our SaaS offering would could be tuned if the permissions stuff
happened on CPU optimized nodes and the API was running on network optimized
nodes. Or maybe your CTO has just been raving about this microservices
hotness and you're feeling the pressure. Its time to break the monolith up.

To start out, lets create a server-client pair implementation of our
permissions service. To create a server we simply wrap another `Service`
implementation up in an `http.Handler`.

```golang
// permissions/remote.go
package permissions

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type server struct {
	Service
}

func NewRemoteServer(s Service) http.Handler {
	return server{s}
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	ok, _ := s.CanI(string(body))
	if ok {
		io.WriteString(w, "true")
	} else {
		io.WriteString(w, "false")
	}
}

type client struct {
	addr string
}

func (c client) CanI(thing string) (bool, error) {
	resp, err := http.Post(c.addr, "text/plain", strings.NewReader(thing))
	if err != nil {
		return false, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "true" {
		return true, nil
	}
	return false, nil
}

func NewRemote(addr string) Service {
	return client{addr}
}
```

With this new networked implementation set up, lets modify our application entry
point to configurate between microservice and monolith operation:

```golang
func main() {
	permsURL := flag.String("perms-url", "", "URL of permissions service")
	flag.Parse()

	switch {
    case len(os.Args) > 1 && os.Args[1] == "permissions":
        // If invoked with `./microlith permissions` we're just running the permissions
        // microservice
		service := permissions.NewSimple()
		handler := permissions.NewRemoteServer(service)
		http.ListenAndServe(":8081", handler)
    default:
        // If invoked as `./microlith` we're running the product main entrypoint
		if *permsURL == "" {
            // No permiossions service URL provided. We're running
            // in monolith mode
			http.ListenAndServe(":8080", API{permissions.NewSimple()})
		} else {
            // Use the permissions URL provided to make a remote client. We're
            // in microservice mode.
			http.ListenAndServe(":8080", API{permissions.NewRemote(*permsURL)})
		}
	}
}
```

Now we can spin up our permissions server by running

```shell
$ microlith permissions
```

and our main entry point with

```shell
$ microlith -perms-url http://localhost:8081
```

## Final Thoughts

Seeing a concrete example, some of the advantages and disadvantages become
clear. If you're disciplined, you'll always have a fully functional version
of your product one command away. The local testing and iteration this is an
obvious win, but there are some subtle wins here as well. What of you notice
a big opportunity to pivot towards an on-prem offering? Or perhaps land a
client that requires an air-gapped deployment? In these spaces, easy of
operations and simple deployment story _is_ the product. Being able to spin
up an all-in-one version puts your application miles ahead of the competition
in these circumstances.

Some pain points still remain obviously, as these solutions scales you still
need to maintain that common entry point because its everyones executeable.
There are probably some opportunities there for muxing commandling args and
sending the details of parsing and launch down into each microservice
library.

Another aspect at play in this pattern is the need to centralize (for the
most part) on one programming language. At the very least, service
implementators need to create a client in the core language and some kind of
wrapper for the server. Its a hurdle, but I think it is a good source of
friction pushing teams to use the same technologies where possible.

> You can find the code for this example at [https://github.com/nsmith5/microlith][2]

[2]: https://github.com/nsmith5/microlith
