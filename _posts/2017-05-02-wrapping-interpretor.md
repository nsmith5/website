---
layout: post
title: "Interacting with external interpretors in Julia"
date: 2017-05-02 14:15:00
categories: julia interpretors
---

basically, 

```julia

# I/O pipes
input = Pipe()
output = Pipe()

# Your favourite interpreter
cmd = `myinterpreter`

# Start up the process
proc = spawn(cmd, (input, output, STDERR))

# Close the unnecessary sides of the pipes
close(input.out)
close(output.in)

# Read and write away!
write(input, "my-expressions;")
out = readavailable(output) |> String
```

Is good