---
layout: post
title:  "Python-style Generators in Julia"
date:   2017-05-02 13:39:00
categories: julia generators
---

Blah blah, generators

```julia
julia> using PyGen

julia> @pygen function pascal(n)
	      i = 0
	      while i <= n
	          yield(i)
	          i += 1
	      end
	  end
pascal (generic function with 1 method)

julia> n = 20;

julia> sum(pascal(n)) == n * (n + 1) / 2
```
