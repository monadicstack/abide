# Bozo

Bozo is a code generator and runtime library that helps you
write RPC-enabled and Event Sourced (micro) services and APIs.
It parses the interfaces/structs/comments in your code service
code to generate all of the client, server, gateway, and pub-sub
communication code automatically.

This is the spiritual successor to [Frodo](https://github.com/monadicstack/frodo).
Bozo supports all of Frodo's RPC/HTTP related features, but it addresses
some shortcomings in the architecture and adds in the ability to
seamlessly have your services support both RPC 
and Event-Driven communication with almost no extra code. The code
generation does all the hard stuff for you.
