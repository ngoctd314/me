---
title: "Channel in Golang"
date: 2022-11-03
authors: ["NgocTD"]
draft: true
tags: ["basic", "golang"]
---

Channel is an important built-in feature in Go. It is one of the features makes Go unique. Channel makes concurrent programming convenient, fun and lowers the difficulties of concurrent programming. Channel mainly acts as a concurrency synchronization technique. To understand channels better, the internal structure of channels and some implementation details by the standard Go compiler/runtime are also simply described.

## Channel Introduction

Don't communicate by sharing memory, share memory by communicating.

Communicating by sharing memory and sharing by communicating are two programming manners in concurrent programming. When goroutines communicate by sharing memory, we use traditional concurrency sychronization techniques, such as mutex locks, to protect the shared memory to prevent data racts.

Go also provides another concurrency sychronization technique, channel. Channels make goroutines share memory by communicating. We can view a channel as an internal FIFO queu within a program. Some goroutines send values to the queue (the channel) and some other goroutines receive values from the queue. 

Along with transfering values (through channels), the ownership of some values may also be transferred between goroutines (ownership on logic view). When a goroutine send a value to a channel, we can view the goroutine releases the ownership of some values. When a goroutine receives a value from a channel, we can view the goroutine acquires the ownerships of some values.

## Channel Value Comparisons

All channel types are comparable types.
If one channel value is assigned to another, the two channels share the same underlying part(s). In other words, those two channels represent the same internal channel object. The result of comparing them is true.

## Detailed Explanations for Channel Operations

|Operation|A Nil Channel|A Closed Channel|A Not-Closed Non-Nil Channel|
|-|-|-|-|
|Close|panic|panic|success|
|Send Value To|block for ever|panic|block or succeed to send|
|Receive value from|block for ever|never block|block or success to receive|

To bettern understand channel types and values, and to make some explainations easier, looking in the raw internal structures of internal channel objects is very helpful.

We can think of each channel consistin of three queues internally:

1. The receiving goroutine queue (generally FIFO). The queue is a linked list without size limitation. Goroutines in this queue are all in blocking state and waiting to receive values from that channel.

2. The sending goroutine queue (generally FIFO). The queue is also a linked list without size limitation. Goroutines in this queue are all in blocking state and waiting to send values to that channel.

3. The value buffer queue (absolutely FIFO). This is a circular queue. Its size is equal to the capacity of the channel. If the current number of values stored in the value buffer queue of the channel reaches the capacity of the channel, the channel is called in full status. If no values are store in the value buffer queue of the channel currently, the channel is called in empty status. For a zero-capacity (unbuffered) channel is also in both full and empty status.

**Each channel internally holds a mutex lock which is used to avoid data races in all kinds of operations**

### Channel operation: try to receive

When a goroutine R tries to receive a value from a not-closed non-nil channel, the goroutine R will acquire the lock associated with the channel firstly, the do the following steps until one condition is satisfied.

### Channel operation: try to send

When a goroutine S tries to send a value to a not-closed non-nil channel, the goroutine S will acquire the lock associated with the channel firstly, then do the following steps until one step condition is satisfied.

### Channel operation: try to close

When a goroutine tries to close a not-closed non-nil channel, once the goroutine has acquired the lock of the channel, both of the following two steps will be performed by the following order.

**After a non-nil channel ise closed, channel receive operations os the channel will never block**

The values in the value buffer of the channel can still be received. The second optional bool return values are still true. Once all the values in the buffer are taken out and received, infinite zero values of the element type of the channel will be received by any of the following        receiving operations on the channel.

**References** [Channel Use Case go101](https://go101.org/article/channel.html)