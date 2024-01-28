# gRPC

## Introduction

gRPC can use protocol buffers as both its Interface Defination Language and as its underlying message interchange format.

In gRPC, aclient application can directly call a method on a server application on a different machine as if it were a local object.

![image](https://grpc.io/img/landing-2.svg)

### Protocol Buffers

By default, gRPC uses Protocol Buffers, a mechanism for serializing structured data. 

- First step when working with protocol buffers is to define the structure for the data you want to serialze in a ***proto file***: an ordinary text file with `.proto` extension. Protocol buffer data is structured as ***messages***, where each message is a small logical record of information containing a serial of name-value paris called ***fields***.

```protobuf
message Person {
  string name = 1;
  int32 id = 2;
  bool has_ponycopter = 3;
}
```

- Then once you've specified your data structures, you use the protocol buffer complier `protoc` to generate data access classes.

```protobuf
// The greeter service defination
service Greeter {
	//Sends a greeting
	rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

gRPC uses `protoc` with a special gRPC plugin to generate code from proto file.

## Service methods

- Unary RPC: send a single request and get a single response.

  ```protobuf
  rpc SayHello(HelloRequest) returns (HelloResponse);
  ```

- Server-side streaming RPC: send a request and get a stream to read a sequence of messages.

  ```protobuf
  rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse);
  ```

- Client-side streaming RPC: streaming RPCs where the client writes a sequence of messages and sends them to the server, wait for the server to read them and return its response. gRPC guarantees message ordering within an individual RPC call.

  ```protobuf
  rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse);
  ```

- Bidiretional streaming RPC: Client and server can read and write in whatever order they like.

  ```protobuf
  rpc BidiHello(stream HelloRequest) returns (stream HelloResponse);
  ```

## Install

### protobuf

```sh
brew install protobuf
```

### Plugin

```sh
brew install protoc-gen-go
brew install protoc-gen-go-grpc
```

### Go protobuf

```sh
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
```

### grpc-go

```sh
go get -u google.golang.org/grpc
```

## Basic

### Define the service

```protobuf
syntax = "proto3";

option go_package = "github.com/User/grpc/hello";
package hello;

service Hello {
    //SayHello method
    rpc SayHello (HelloRequest) returns (HelloReply) {}

    //LotsOfReplies method
    rpc LotsOfReplies (HelloRequest) returns (stream HelloReply) {}
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

Specifying their request and response types within the four service methods.

### Generating client and server code

Open the directory where `.proto` lies.

```sh
~/ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./hello.proto
```

- `hello.pb.go`, which contains all the protocol buffer code and response message types.
- `hello_grpc.pb.go`
  - An interface type(or *stub*) for clients to call with the methods defined in the `Hello` service.
  - An interface type for servers to implement.

### Creating the server

- Implementing the service interface generated from our service definition.
- Running a gRPC server to listen for requests from clients and dispatch them to the right service implementation.

