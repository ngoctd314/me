---
title: "Channel in Golang"
date: 2022-11-05
authors: ["ngoctd"]
draft: false
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

1. Check buffer, if the value buffer queue of the channel is not empty. The receiving goroutine queue of the channel must be empty ( buffer != empty => receiveing queue == emtpy ). The goroutine R will receive (by unshifting) a value from the value buffer queue. If the sending goroutine queue of the channel is also not empty, a sending goroutine will be unshifted out of the sending goroutine queue and resumed to running state again. The value the just unshifted sending goroutine trying to send will be pushed into the value buffer queue of the channel. The receiving goroutine R continues running. For this scenario, the channel receive operation is called a non-blocking operation.

- The goroutine R will receive a value from the value buffer queue.

![receive value from buffer, sending goroutine queue is empty](../../images/channel%20golang/receive1.png)

- The goroutine R will recceive a value from the value buffer queue. Sending goroutine is not empty. Goroutine S send value to buffer and enter running state again.

![receive value from buffer, sending goroutine is not emtpy](../../images/channel%20golang/receive2.png)

**-> The receiving goroutine R continues running. The channel receive operation is called a non-blocking operation**

2. Check buffer, the value buffer of the channel is empty. If the sending goroutine queue of the channel is not empty, in which case the channel must be an unbuffered channel, the receiving goroutine R will unshift value from a send goroutine. The just unshifted sending goroutine will get unblocked and resumed to running state again. 

![receive value from buffer, sending goroutine is not emtpy](../../images/channel%20golang/receive3.png)

**-> The receiving goroutine R continues running. The channel receive operation is called a non-blocking operation**

3. If value buffer queue and the sending goroutine queue of the channel are both emtpy, the goroutine R will be pushed into the receiving goroutine queue of the channel and enter (and stay in) blocking state. It may be resumed to running state when another goroutine sends a value to the channel later.

![receive value from buffer, sending goroutine is not emtpy](../../images/channel%20golang/receive4.png)

**-> The receiving goroutine R enter blocking state. The channel receive operation is called a blocking operation**

### Channel operation: try to send

When a goroutine S tries to send a value to a not-closed non-nil channel, the goroutine S will acquire the lock associated with the channel firstly, then do the following steps until one step condition is satisfied.


1. Check receiving goroutine queue. If the receiving goroutine queue of the channel is not empty, in which case the value buffer queue of the channel must be empty, the sending goroutine S will unshift a receiving goroutine from the receiving goroutine queue of the  channel and send the value to the just unshifted receiving goroutine. The just unshifted receiving goroutine will get unblocked and resumed to running state again.

![send value](../../images/channel%20golang/send1.png)

**-> The sending goroutine S continues running. The channel send operation is called a non-blocking operation**

2. Check receiving goroutine queue (empty), check buffer queue ( not full ), in which case the sending goroutine queue must be also empty, the value the sending goroutine S trying to send will be pushed into the value buffer queue.

![send value](../../images/channel%20golang/send2.png)

**-> The sending goroutine S continues running. The channel send operation is called a non-blocking operation**

3. Check receiving goroutine queue (empty), check buffer queue ( full ), the sending goroutine S will be pushed into the sending goroutine queue of the channel and enter (and stay in) blocking state. It may be resumed to running state when another goroutine receives a value from the channel later.

![send value](../../images/channel%20golang/send3.png)

**-> The sending goroutine S enter blocking. The channel send operation is called a blocking operation**

Once a non-nil channel is closed, sending a value to the channel will produce a runtime panic in the current goroutine. Note sending data to a closed channel is viewed as a non-blocking operation.

### Channel operation: try to close

When a goroutine tries to close a not-closed non-nil channel, once the goroutine has acquired the lock of the channel, both of the following two steps will be performed by the following order.

1. If the receiving goroutine queue of the channel is not empty, in which case the value buffer of the channel must be empty, all the goroutines in the receiving goroutine queue of the channel will be unshifted one by one, each of themm will receive a zero value of the elemenet type of the channel and be resumed to running state.

![send value](../../images/channel%20golang/close1.png)

2. If the sending goroutine queue of the channel is not empty, all the goroutines in the sending goroutine queue of the channel will be unshifted one by one and each of them will produce a panic for sending on a closed channel. This is the reason why we should avoid concurrent send and close operations on the same channel.

After a channel is closed, the values which have been already pushed into the value buffer of the channel are still there.

**After a non-nil channel is closed, channel receive operations os the channel will never block**

### Some facts about the internal queues of a channel

- If the channel is closed, both its sending and receiving goroutine queue must be empty, but its value buffer may not be empty.
- At any time, if the value buffer is not empty, then its receiving goroutine queue must be empty.
- At any time, if the value buffer is not full, then its sending goroutine queue must be empty.
- If the channel is buffered, then at time, at least one of the channel's goroutine queues must be empty (sending, receiving or both).
- If the channel is unbuffered, most of the time one of its sending goroutine queue and the receiving goroutine queue must be empty, with one exception. The exception is that a goroutine may be pushed into both of the two queues when execution a select control flow code block.

**References** [Channel Use Case go101](https://go101.org/article/channel.html)