---
title: "Designing An Augmented Backus-Naur Form Parser Generator"
date: 2020-09-14T19:57:30-07:00
draft: false
---

I've very interested in email. Email is an important piece of establishing a
digital identity. Email is used, either directly or indirectly, to establish
your identity with almost every single service on the internet. Unfortunately
email is extremely difficult to self host. I've tried several times and at this
point I've given up and host my email using GSuite (If you can't beat them,
join them I guess). As a result, I've dreamt on and off about creating my own
email server. The challenge is that email is specified by a handful of
protocols that, in turn, are specified with a spawling mountain of
specifications.

Augmented Backus-Naur Form (ABNF) is many of these specifications. For
instance, it specifies things like:

- _What is the allowed format of Email messages?_ | [RFC 5532][RFC5532]
- _What are allowed commands in the Simple Mail Transfer Protocol?_ | [RFC
  5321][RFC5321]

[RFC5532]: https://tools.ietf.org/html/rfc5322
[RFC5321]: https://tools.ietf.org/html/rfc5321

For example, the `EHLO` command SMTP clients use is specified with the section
of ABNF below. This and other rules combine to complete a complete picture of
everything that is allowed in the context of the SMTP protocol 

```abnf
ehlo-ok-rsp    = ( "250" SP Domain [ SP ehlo-greet ] CRLF )
                    / ( "250-" Domain [ SP ehlo-greet ] CRLF
                    *( "250-" ehlo-line CRLF )
                    "250" SP ehlo-line CRLF )
ehlo-greet     = 1*(%d0-9 / %d11-12 / %d14-127)
               ; string of any characters other than CR or LF
ehlo-line      = ehlo-keyword *( SP ehlo-param )
ehlo-keyword   = (ALPHA / DIGIT) *(ALPHA / DIGIT / "-")
               ; additional syntax of ehlo-params depends on
               ; ehlo-keyword
ehlo-param     = 1*(%d33-126)
               ; any CHAR excluding <SP> and all
               ; control characters (US-ASCII 0-31 and 127
               ; inclusive)
```

Ideally, I'd like to just copy and paste these specifications into the code for
my email server so that I know I'm 100% compliant with the RFC (At least as far
as the grammar is concerned). What might that look like? Perhaps something like

```golang
package main

import "github.com/nsmith5/abnf"

const grammar = `
ehlo-ok-rsp    = ( "250" SP Domain [ SP ehlo-greet ] CRLF )
                    / ( "250-" Domain [ SP ehlo-greet ] CRLF
                    *( "250-" ehlo-line CRLF )
                    "250" SP ehlo-line CRLF )
ehlo-greet     = 1*(%d0-9 / %d11-12 / %d14-127)
ehlo-line      = ehlo-keyword *( SP ehlo-param )
ehlo-keyword   = (ALPHA / DIGIT) *(ALPHA / DIGIT / "-")
ehlo-param     = 1*(%d33-126)
`

func main() {
    parser, err := abnf.New(grammar)
    if err != nil {
        // Invalid grammar
        panic(err)
    }
    
    parseTree, err := parser.Parse("250 nfsmith.ca\r\n")
    if err != nil {
        // Invalid command
        panic(err) 
    }

    // Use parse tree to do stuff
}
```

Yup. That would definitely be helpful. What do we want from that parse tree
data structure?

```golang
type ParseNode struct {
    Production string
    Matched string
    Children []ParseNode
}
```

That feels useful, because we'd get something like the following for the example above:

```golang
ParseNode{
    Production: "ehlo-ok-resp",
    Matched: "250 nfsmith.ca\r\n",
    Children: []ParseNode{
        ParseNode {
            Production: `"250"`,
            Matched: "250",
            Children: nil, 
        },
        ParseNode{
            Production: `SP`,
            Matched: " ",
            Children: nil,
        },
        ParseNode{
            Production: `Domain`,
            Matched: "nfsmith.ca",
            Children: nil,
        },
    },
}
```

This is fairly awesome because we can extract the information we want into
domain specific data structures afterwards. For instance, we can take the parse
tree above and marshal data into something like a `EhloOkResponse` struct like,

```golang
type EhloOkResponse struct {
    Domain string
}
```

Overall, the component diagram would look like this,

<img src="/img/smtp-diagram.jpg" alt="email server diagram" width=100%></img>

## Follow up

Time to build! If you'd like to follow along, I'll be building at
https://github.com/nsmith5/abnf (NB: the existing implementation in that
repository doesn't implement this design and has many short comings). Working
towards implementing this design and hopefully writing about the experience
over the next little while.
