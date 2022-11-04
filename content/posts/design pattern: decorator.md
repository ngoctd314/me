---
title: "Decorator pattern"
date: 2022-11-04
authors: ["NgocTD"]
draft: true
tags: ["design pattern"]
series: ["head first design pattern"]
---

I used to think real men subclassed everything. That was until I learned the power of extension at runtime, rather than at compile time.

Just call this chapter "Design Eye for Inheritance Guy." We'll re-examine the typical overuse of inheritance and you'll learn how to decorate your classes at runtime using a form of object composition. Why? Once you know the techniques of decorating, you'll be able to give your objects new responsibilities without making any code changes to the underlying classes.

## Welcome to Starbuzz Coffee

Starbuzz Coffee has made a name for itself as the fastest-growing coffee shop. When they first went into business they designed their classes like this...

![inheritance version](../../images/design-patterns/decorator/inheritance.png)

## Summary