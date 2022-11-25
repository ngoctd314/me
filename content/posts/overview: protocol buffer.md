---
title: "Protocol Buffers Overview"
date: 2022-11-07
draft: true
tags: ["grpc"]
authors: ["ngoctd"]
---

Protocol buffers provide a language-neutral, platform-neutral, extensible mechanism for serializing structured data in a forward-compatible and backward-compatible way. It's like JSON, except it's smaller and faster, and it generates native language bindings.

Protocol buffers are a combination of the definition language (created in .proto files), the code that the proto compiler generates to interface with data, language-specific runtime libraries, and the serialization format for data that is written to a file (or sent across a network connection).

**Some of the advantages of using protocol buffers**

- Compact data storage
- Fast parsing
- Availability in many programming languages
- Optimized functionality through automatically-generated classes

**When are Protocol Buffers not a Good Fit?**