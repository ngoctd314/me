---
title: "wire for D. D for Dependency Injection"
date: 2022-12-12
draft: false
authors: ["ngoctd"]
tags: ["golang", "basic", "solid"]
---

Wire has two core concepts: providers and injectors.

## Providers

The primary mechanism in Wire is the **providers**: a function that can produce a value.

```go
type Foo struct {
    X int
}

func ProvideFoo() Foo {
    return Foo{X: 42}
}
```

Providers can also return errors

```go

```
## Injectors

An application wires up these providers with an **injector**: a function that calls providers in dependency order. With Wire, you write the injector's signature, then Wire generates the function's body.