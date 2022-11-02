---
title: "Golang All in One"
date: 2022-11-02T07:33:22+07:00
draft: true
tags: ["golang", "basic"]
categories: ["cat1"]
---

All about golang at basic level in one post. Content in this post i have learned in [go101 book ](https://go101.org/) and on the internet. At this time, i think it can be misconcept in golang.


## Basic Types and Basic Value Literals

rune type (a.k.a int32 type), are special integer types. A rune value is intended to store a Unicode point. A rune literal is expressed as one or more characters enclosed in a pair of quotes, for example 'a' (the Unicode value of character a is 97). We should know that some Unicode characters are composed of more than one code points, for example 'a', '\x61', '\141' are the same (the Unicode value is 97).

```go
func main() {
	println('a' == 97)
	println('a' == '\x61')
	println('a' == '\141')
	println('a' == '\u0061')
}
```