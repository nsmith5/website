---
title: "Taking the bite out of x509 certificates with the step CLI"
date: 2022-01-16T12:22:32-08:00
draft: false
---

Parsing, generating and troubleshooting certificates is critical skill in
developing web services. Certificates establish trust on the web (e.g that
indeed the company Google is serving you content when you go to www.google.com)
and to encrypt traffic once trust is established using TLS.

The defacto developer tooling in this space is `openssl`. `openssl` is a
cryptograph swiss army knife. If you need to get something done in the broad
space of cryptography, `openssl` _can_ do it. The only problem is figuring
out the right magical incantation to make it happen.

I, like many folks I imagine, have an every growing personal wiki page of how
to get stuff done with openssl. Here are some of the highlights from mine:

```bash
# Connect to a client using TLS (dumps a mountain of information on 
# the entire handshake process)
openssl s_client -connect example.com:443

# Parse x509 certificate on disk
cat cert.pem | openssl x509 -text -noout

# Parse x509 certificate from a server
echo | \
  openssl s_client -servername example.com -connect example.com:443 2>/dev/null |\
  openssl x509 -text -noout

# Create a self-signed certificate
openssl req -x509 \
  -newkey ed25519 \
  -sha256 \
  -keyout key.pem \
  -out cert.pem \
  -subj "/CN=example.com" \
  -days 365
```

Unfortunately, sometimes you need to do something with `openssl` _that you know
it can do_ but you simply can't find out how to do. I hit this wall recently
while trying to create some test data for an open source project I was working
on. I needed to create a root certificate authority and a chain of at least 2
intermediate certificates authorities from the root.

To be clear, I know in my heart that `openssl` can do this. The challenging
thing with `openssl` is just getting all the flags and things just right.
Sometimes you'll get lucky with your search engine-fu and find your use case
perfectly documented, but alas it didn't get so lucky.

## Enter the step CLI

The `step` CLI is [newer crypto swiss army
knife](https://smallstep.com/blog/zero-trust-swiss-army-knife/). I'd heard of
the `step` CLI when I read [this excellent deep
dive](https://smallstep.com/blog/everything-pki/) on public key infrastructure
from one of the `step` authors. Like that deep dive article, the `step` CLI
has a strong focus on modern standards. One example to highlight this: By default
`openssl` is going to prompt you for "distinguished names" when you generate
an x509 certificate. You may find yourself contemplating the _Organizational
Unit_ of your server and scratch your head wondering why things got so
corporate all of a sudden. With `step` you'll have to go out of your way to
specify this information. This is a good thing because this silly distinguished
names have have been deprecated by the [CAB
forum](https://cabforum.org/wp-content/uploads/CA-Browser-Forum-BR-1.6.1.pdf)
(ref section 7.1.4.2) and made optional.

Alright so good tools make the easy stuff easy right? Lets take a look at
how `step` makes the easy use cases easy. A simple self-signed certificate:

```bash
# Generate self-signed certificate
step certificate create --profile self-signed --subtle subject cert.pem key.pem

# Inspect certificate
step certificate inspect cert.pem

# Create root CA and intermediate CA
step certificate create --profile root-ca root root.cert.pem root.key.pem
step certificate create --profile intermedate-ca \
  --ca root.cert.pem \
  --ca-key root.key.pem \
  intermediate intermediate.cert.pem intermediate.key.pem

# Sign a leaf certificate with the intermediate above
step certificate create --profile leaf \
  --ca intermediate.cert.pem \
  --ca-key intermediate.key.pem \
  my-leaf leaf.cert.pem leaf.key.pem
```

Pretty straight forward right? The `profile` flag is really helpful for setting
a lot of the right defaults for common use cases.

## Templating

The concern with a simplistic interface like the one before is that we
might have lost some control and that the system might not work for
a particular corner case. Indeed my original test data use case requires

- 1 root CA
- 2 intermediates CA
- Each intermediate CA needs to have the code signing extended key usage flag set

This is because the application I was working on [must issue code signing
certificates](https://github.com/sigstore/fulcio/) and in accordance with the
CAB forum we need to have all intermediates share the extended key usages (EKU)
that they will issue to leaf certificates. This is called EKU chaining.

This is where we start to use the templating system of `step`. In fact, the
profiles used above were simply preset templates. To start we need to create a
root CA with max path length = 2. This ensures that the number of hops to the
leaf certificate is at most 3 (yeah weird, but its one of those off by one
moments you know and love from computer science). This is achieved with the
following template and `step` invocation:

```bash
cat <<EOF > root.tpl
{
   "subject": {
      "commonName": "root"
   },
   "issuer": {
      "commonName": "root"
   },
   "keyUsage": ["certSign", "crlSign"],
   "basicConstraints": {
     "isCA": true,
     "maxPathLen": 2
   }
}
EOF

step certificate create --template root.tpl root root.cert.pem root.key.pem
```

Now we create the intermediates, again using the templating system to control EKU
as desired:

```bash
cat <<EOF > intermediate1.tpl
{
   "subject": {
      "commonName": "intermediate1"
   },
   "keyUsage": ["certSign", "crlSign"],
   "extKeyUsage": ["codeSigning"],
   "basicConstraints": {
     "isCA": true,
     "maxPathLen": 1
   }
}
EOF

step certificate create \
  --template intermediate1.tpl \
  --ca root.cert.pem \
  --ca-key root.key.pem \
  intermediate1 intermediate1.cert.pem intermediate1.key.pem

cat <<EOF > intermediate2.tpl
{
   "subject": {
      "commonName": "intermediate2"
   },
   "keyUsage": ["certSign", "crlSign"],
   "extKeyUsage": ["codeSigning"],
   "basicConstraints": {
     "isCA": true,
   }
}
EOF

step certificate create \
  --template intermediate2.tpl \
  --ca intermediate1.cert.pem \
  --ca-key intermedate1.key.pem \
  intermediate2 intermediate2.cert.pem intermediate2.key.pem
```

Hazzah! Not too bad at all. I think the key thing about `step` is that, not
only was it easy to figure this out once I understood the basics, the
documentation had a [massive amount of
examples](https://smallstep.com/docs/step-cli/reference/certificate/create#examples)
to tweak and extend.

## Conclusions

If you're looking for a well documented, versitile tool crypto tool without all the
sharp edges: the `step` CLI is for you. This is especially true if you haven't already
invested a load of time learning about `openssl`.

- You'll figure out how to achieve your goals faster because of the quality
  docs
- Modern standards are the defaults in `step` (it's harder (maybe impossible?)
  to accidentally use SHA1 for example)
