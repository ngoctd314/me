---
title: "GRPC Overview"
date: 2022-11-07
draft: true
tags: ["grpc"]
authors: ["ngoctd"]
---

gRPC is a modern, open source remote procedure call (RPC) framework that can run anywhere. It enables client and server applications to communicate transparently, and makes it easier to build connected systems.

In gRPC, a client application can directly call a method on a server application on a different machine as if it were a local object, making it easier for you to create disitributed applications can services. As in many RCP systems, gRPC is based around the idea of defining a service, specifying the methods that can be called remotely with their parameters and return types. On the server side, the server implements this interface and runs a gRPC server to handle client calls. On the client side, the client has a stub (referred to as just a client in some languages) that provides the same methods as the server.

## Service definition

Like many RPC systems, gRPC is based around the idea of defining a service, specifying the methods that can be called remotely with their parameters and return types.

```proto
service HelloService {
    rpc SayHello(HelloRequest) returns (HelloResponse);
}

message HelloRequest {
    string greeting = 1;
}

message HelloResponse {
    string reply = 1;
}
```

## Four kinds of service method

- Unary RPCs where the client sends a single request to the server and gets a single response back

```proto
service UnaryService {
    rpc Unary(Requestr) returns (Response);
}
```

- Server streaming RPCs where the client sends a request to the server and gets a stream to read a sequence of messages back. The client reads from the returned stream until there are no more messages. gRPC guarantees message ordering within an individual RPC call.

```proto
service ServerStreaming {
    rpc Streaming(Request) returns (Responses);
}
```

- Client streaming RPCs where the client writes a sequence of messages and sends them to the server, again using a provided stream. Once the client has finished writing the messages, it waits for the server to read them and return its response. Again gRPC guaranteed message ordering within an individual RPC call.

```proto
service ClientStreaming {
    rpc Streaming(Requests) returns (Response);
}
```

- Bidirectional streaming RPCs where both sides send a sequence of messages using a read-write stream. Server could wait to receive all the client messages before writing its responses, or it could alternately read a message then write a message, or some other combination of reads and writes. The order of messages in each stream is preserved.

```proto
service BidirectionalStreaming {
    rpc Streaming(Requests) returns (Responses)
}
```

## Using the API

gRPC provides protocol buffer compiler plugins that generate client and server side code. gRPC users typically call these APIs on the client side and implement the corresponding API on the server side.

- On the server side, the server implements the method declared by the service and runs a gRPC server to handle client calls. The gRPC infrastructure decodes incoming requests, executes service methods, and encodes service response.
- On the client side, the client has a local object known as stub that implements the same methods as the service. The client can then just call those methods on the local object, wrapping the parameters for the call in the appropriate protocol buffer message type.

## RPC life cycle

What happen when a gRPC client calls a gRPC server method.

### Unary RPC

Client sends a single request and gets back a single response.


### Server streaming RPC

### Client streaming RPC

### Bidirectional streaming RPC

## Deadlines/Timeouts

gRPC allows clients to specify how long they are willing to wait for an RPC to complete before the RPC is terminated with a DEADLINE_EXCEEDED error.

## RPC termination

In gRPC, both the client and server independent and local determinations of the success of the call, and their conclusions may not match. This means that, for example, you could have an RPC that finishes successfully on the server side ("I have sent all my responses!") but fails on the client side ("The responses arrived after my deadline!"). It's also possible for a server to decide to complete before a client has sent all its requests.

## Cancelling an RPC

Either the client or the server can cancel an RPC at any time. A cancellation terminations the RPC immediately so that no further work is done.

## Channels

A gRPC channel provides a connection to a gRPC server on a specified host and port. A channel has state, including connected and idle.