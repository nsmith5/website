---
title: "One Month in Open Source"
date: 2022-02-10T09:35:49-08:00
draft: false
---

Roughly one month ago, I left my job to work in open source. This is second
time open source had a hand in me quitting something big and I'm hoping that
this time is the last.

## Leaving academia

The first time was in quitting my PhD program. Two years into a PhD in classical
condensed matter physics, I was doing a lot of simulation work. My work would
start with some pen and paper theory and then be translated into numerical
simulation. These simulations were distributed programs that would run on the
university's super computer. The output was often quite beautiful when rendered
into [movies or images](https://www.physics.mcgill.ca/~provatas/.index.html) so
we lovingly called our graduate student office the "pretty picture department".
Academia, and physics in particular, has a long history of creating and using
open source software and my department was no exception. This is where I got my
first introduction to Linux and where I learned to cobble together open source
libraries like [FFTW](https://www.fftw.org/) and
[OpenMPI](https://www.open-mpi.org/) to create my simulation code.

While all of our code was built on open source software, there was a stark
contrast between the collaborative model used to build these tools and how
research itself proceeded. While our research findings were openly published,
the academic community is more competitive than collaborative. Simulation code
is rarely published so its difficult to build directly upon each others work.
Even simply reproducing a paper from outside of your own research group can be
extremely difficult.

It was natural to contrast this with my experience submitting [my first small
open source patches](https://github.com/JuliaLang/julia/pull/21626) to projects
like the [Julia programming language](https://github.com/JuliaLang/julia). This
was just a documentation change, but review was thorough and everyone was
really encouraging as I fumbled my way a long.

In short, open source software showed me:

- Open collaboration on technical problems was possible
- I could access some of the best minds in the field by volunteering my time on
  these problems
- Folks were willing to teach me and help me join this community in a
  productive way

I cut my PhD short into a Masters program and hopped into the first software
role I could find.

## Leaving proprietary software 

Fast-forward 5 years and I've been in various software roles from data science
to infrastructure. Contributing to open source software and keeping up with
what the community is building has been my constant hobby and companion in that
journey and its been a huge asset in my career.

Building in open source has helped me:

- Discover great models of collaboration to bring into the workplace,
- Know how to create a great pull request (context, context, context!), and
- When something is already a solved problem ("Hey, did you know there is
  already an open specification for X?")

That said, it always remained a hobby or a minor role in my work. I'd submit
little bug fixes here and there or deploy an open source project at to solve
some problem, but it was always outside the main focus.

I was always a little jealous of those folks that were constantly working in
open source as part of their career, but it can be easy to write yourself out
of ever having this career yourself. The most visible folks being paid to work
in these spaces are often extremely senior engineers from prestigious
companies. Want to contribute to Kubernetes? You shouldn't be surprised if your
code ends up being reviewed by a clever Google engineer for example. While this
means you get to learn a lot by contributing, it can make it hard to imagine
yourself on the other side of that review.

That said, if you never fully engage in those communities how could you ever
really hope to find yourself being paid to contribute to them? If it's your
after work hobby it's getting the fumes of your empty tank of motivation you
know? 

## One month in open source

So the idea was this: just dive in there full time and give it your full
attention. Build things, meet folks and contribute full time. After giving
myself at least a month I'd see if I could land a role that would let me stay
active in open source.

It's important to mention at this point the absolutely massive pile of luck and
privilege required to do something like this. Its truly rare to simply be able
to fuck off from your career for a while to try something like this out. In my
case it required knowing I could land a job again whenever I wanted it, having
the money to live without work for a long time and a partner that actively
pushed me towards the idea because I was excited about it.

I chose the [Sigstore](https://www.sigstore.dev/) project to focus my effort
and in particular the [Fulcio](https://github.com/sigstore/fulcio/) certificate
authority. I structured my time contributing to the project much like I would a
full-time job. I spent about 60% of my time actively trying to fix issues and
push the project along and then about 40% of my time trying to "onboard". For
the Sigstore project this mean playing with and learning about software supply
chain security how ever I could.

Being open about what I was doing with the community meant a lot of people
helped out. I openly stated I was looking to push my career towards open source
supply chain security and was looking for mentorship and guidance if anyone was
willing to help out. Many folks offered support and let me collaborate on their
work. There are many many ways that folks helped out with this. Here are a
couple examples that meant a lot to me:

- Matt Moore and Ville Aikas immediately shared some work on using file based
certificate backends for Fulcio I could help out with
- Scott Nichols taught me all about [Cloud events](https://cloudevents.io/) and
helped me wire them up to a project I had created
- Luke Hinds shared his thorough tutorial of setting up sigstore from scratch
- Dan Lorenc, Bob Callaway and Luke Hinds all patiently let me practice using
  the security disclosure for Sigstore because I wanted to learn about it.
- Hayden Blauzvern graciously spent a lot of time combing through my pull
  requests and pointing out references where I could learn more when I did
crimes against x509 repeatedly

Basically, the whole community was immediately receptive and I learned a huge
amount!

## Success?

Alright so, big success right? Yeah absolutely. After a couple weeks I knew
this kind of work was a really great fit for me and I'd met a bunch of
absolutely wonderful folks in the space. Not only that, but I was getting
referrals and interviewing specifically for open source roles. Only two weeks
into the month I was interviewing for three different roles. All three were in
open source and two would keep me working in the Sigstore project in some way!
This was mind blowing to me.

There were aspects of the job interview process that had completely changed as
well. For the roles around Sigstore, I was already working with my potential
future colleges. The risk for them was low and the risk for me was low. I knew
what working together would look like, because we were already doing it. Even
for the role that wasn't connected to Sigstore, when asked "tell me about a
recent technical challenge your had to work through" I could just point to the
sigstore work. It was all open and documented! I just had to walk the
interviewer through the background context.

## Success!!

Ok, so where did I go in the end? I'm so excited that I'll be starting work at
[Chainguard](https://chainguard.dev) this Monday! They're team is full of folks
that have been building awesome tools in the open for years. Tools like:

- [Tekton](https://tekton.dev/)
- [KNative](https://knative.dev/)
- [Minikube](https://minikube.sigs.k8s.io/docs/)
- [Distroless](https://github.com/GoogleContainerTools/distroless)
- [Skaffold](https://skaffold.dev/)

Chainguard is deeply invested in the Sigstore project so I'm extremely excited
that I'll get to continue contributing upstream when I join their team. I'm
completely humbled to be surrounded by folks with so much experience building
in open source and can't wait to help out with their big mission:

> Make software supply chains secure by default ❤️

## This could have been easier

Remember that bit about the how much luck and privilege I needed to take the
time to do this? In retrospect, the leap of faith was absolutely worth it and I
really wish I'd done something similar earlier, but the barrier to even thinking
about doing something like this was really high!

Like I said, I had financial and emotional support in taking this risk, but I
also already knew the right communities to participate in and already know a
lot about open source contribution. It made my time very effective and that is
a lot of what lead to my success.

This could have been easier though and I'd like to make it easier for folks
that want to take the same leap. If you're interested in get into open source
and want some mentorship or advice about how to get started please reach out!
I'd love to help others see the impact they can make and the communities that
they can be a part of ❤️
