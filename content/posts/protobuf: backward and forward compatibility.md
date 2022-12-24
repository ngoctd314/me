## What exactly are protocol buffers?

Protocol Buffers provide a language-neutral, platform-neutral, extensible mechanism for serializing structured data in a forward-compatible and backward-compatible way. It's like JSON, except it's smaller and faster, and it generates native language bindings.

## Some of the advantages of using protocol buffers include:

- Compact data storage
- Fast parsing
- Availability in many programming languages
- Optimized functionality through automatically-generated classes

## Understanding Backward and Forward Compatibility

There are at least two parts to any protobuf system, the sender and the receiver. If either one can be upgraded to a new message format, and the system functionality continues uninterrupted then the message protocol is both forward and backward compatible.

**Backward Compatibility**

If a client that was updated to a new message type but it still able to understand the previous message type then the message change is backward compatible. Backward compatibility is being able to understand messages from a previous version.

**Forward Compatibility**

If a message if changed and a non-updated client can still understand and process the message then the message change if forward compatible. Forward compatibility is being able to understand messages from a future version.

With Protocol buffers, it a sender is upgraded, the receiver can still understand messages if it is forward compatible.