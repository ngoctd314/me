---
title: "12 Factor App"
date: 2023-04-12T09:15:34+07:00
draft: true
authors: ["ngoctd"]
categories: ["architecture"]
---

The Twelve-Factor App methodology is a methodology for building software-as-a-service applications by Adam Wiggins.

## Codebase

One codebase tracked in revision control, many deploys

If there are multiple codebase, it's not an app - it's a distributed system. Each component in a distributed system is an app, and each can individually comply with twelve-factor.

In a microservices-focused world, all the overhead needed to maintain multiple repositories can be tiresome.

## Dependencies

## Config

## Backing services

## Build, release, run

## Processes

## Port binding

## Concurrency

## Disposability

## Dev/pod parity

## Logs

## Admin processes