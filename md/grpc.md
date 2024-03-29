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

option go_package = "github.com/Ferriem/grpc/code/HelloWorld/hello";
package hello;

service Hello {
    //SayHello method
    rpc SayHello (HelloRequest) returns (HelloReply) {}

    //LotsOfReplies method
    rpc LotsOfReplies (HelloRequest) returns (stream HelloReply) {}

    //LotsOfGreetings method
    rpc LotsOfGreetings (stream HelloRequest) returns (HelloReply) {}

    //BidiHello method
    rpc BidiHello (stream HelloRequest) returns (stream HelloReply) {}
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

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/Ferriem/grpc/code/HelloWorld/hello"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Hello_LotsOfRepliesServer) error {
	log.Printf("Received: %v", in.GetName())
	for i := 0; i < 10; i++ {
		stream.Send(&pb.HelloReply{Message: "Hello " + in.GetName() + fmt.Sprintf(" %d", i)})
	}
	return nil
}

func (s *server) LotsOfGreetings(stream pb.Hello_LotsOfGreetingsServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "Hello " + in.GetName()})
		}
		if err != nil {
			return err
		}
		log.Printf("Received: %v", in.GetName())
	}
}

func (s *server) BidiHello(stream pb.Hello_BidiHelloServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Received: %v", in.GetName())
		stream.Send(&pb.HelloReply{Message: "Hello " + in.GetName()})
	}
}

func main() {
	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &server{})
	log.Printf("Starting server on port %d", *port)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

```

- `port = flag.Int("port", 50051, "The server port")` set the port with default value **50051**, or we can specify the port by `./program -port 8080`

  `flag.Int` takes three arguments.

  - The name of the flag. ("**port**")
  - The defalue value ('**50051**')
  - A description of the flag ("**The server port**")

  ```sh
  ~/ ./program -h
  	-port int
          The server port (default 50051)
  ```

- `type server struct`: A struct that implements the methods defined in your gRPC service.

  - `pb.UnimplementedHelloServer` means if you haven't implemented particular method in `server`. gRPC will automatically use the default behavior provided by `UnimplementedHelloServer`
  - `SayHello` and `LotsOfReplies` are particular methods.

- `pb.Register[ServiceName]Server(s, &server{})`: register the service implementatio with the gRPC server.

- `s.Serve()`

  - **Blocking Call**: it will not return until an error occurs or until the server is stopped.

  - **Listening for Requests**: continuously listen for incoming gRPC requests on the specified listener.

  - **Handling Requests**: When a request comes in, the server will handle it by invoking the appropriate gRPC service method based on the RPC requested by the client.
  
- `Server-Side method`: 

  - **Unary RPC** `func (s *MyService) MyMethod(ctx context.Context, request *MyRequest) (*MyResponse, error)`

    - `ctx`  containing deadline, cancellation signals, etc.
    - `request` the signle request parameter sent by the client

  - **Server-Side Streaming RPC** `func (s *MyService) MyMethod(request *MyRequest, stream MyService_MyMethodServer) error`

    - `request`
    - `stream` the server-streaming RPC is represented by a stream that the server can **write** multiple responses to.

    ```go
    type <ServiceName>_FooServer interface {
    	Send(*MsgA) error
    	grpc.ServerStream
    }
    ```

  - **Client-Side Streaming RPC**  `func (s *MyService) MyMethod(stream MyService_MyMethodServer) error`

    - `stream` the server-streaming RPC is represented by a stream that the server can **read** multiple responses to.

    ```go
    type <ServiceName>_FooServer interface {
    	SendAndClose(*MsgA) error
    	Recv() (*MsgB, error)
    	grpc.ServerStream
    }
    ```

  - **Bidirectional Streaming RPC** `func (s *MyService) MyMethod(stream MyService_MyMethodServer) error`

    - `stream` both read and write

    ```go
    type <ServiceName>_FooServer interface {
    	Send(*MsgA) error
    	Recv() (*MsgB, error)
    	grpc.ServerStream
    }
    ```



### Creating the client

```go
package main

import (
	"context"
	"flag"
	"io"
	"log"

	pb "github.com/Ferriem/grpc/code/HelloWorld/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "user"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "the name to hello")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial server:", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)

	res, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: *name})

	if err != nil {
		log.Fatal("Failed to say hello:", err)
	}

	log.Printf("SayHello: %s", res.GetMessage())

	stream, err := c.LotsOfReplies(context.Background(), &pb.HelloRequest{Name: "ferriem"})
	if err != nil {
		log.Fatal("Failed to say hello:", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Failed to recv:", err)
		}
		log.Printf("LotsOfReplies: %s", res.GetMessage())
	}

	var names = []string{"greet1", "greet2", "greet3", "greet4", "greet5"}
	stream2, err := c.LotsOfGreetings(context.Background())
	if err != nil {
		log.Fatal("Failed to say hello:", err)
	}
	for _, name := range names {
		stream2.Send(&pb.HelloRequest{Name: name})
	}
	reply, err := stream2.CloseAndRecv()
	if err != nil {
		log.Fatal("Failed to say hello:", err)
	}
	log.Printf("LotsOfGreetings: %s", reply.GetMessage())

	var notes = []string{"note1", "note2", "note3", "note4", "note5"}
	stream3, err := c.BidiHello(context.Background())
	if err != nil {
		log.Fatal("Failed to say hello:", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream3.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatal("Failed to recv:", err)
			}
			log.Printf("BidiHello: %s", in.GetMessage())
		}
	}()
	for _, note := range notes {
		stream3.Send(&pb.HelloRequest{Name: note})
	}
	stream3.CloseSend()
	<-waitc

}

```

- `grpc.Dial`: create a gRPC channel to communicate with the server. 

- `pb.New[Service name]Client`: Once the gRPC channel is setup, we need a client stub to perform RPCs.

- `stream.Recv()` atomically sort the coming message.

- `Client-Side Method`

  - **Unary Methods** `func YourMethod(ctx context.Context, in *YourRequest) (*YourResponse, error)`

  - **Server-Side Streaming RPC** `func YourMethod(ctx context.Context, in *YourRequest, opts ...grpc.CallOption) (YourService_YourMethodClient, error)`

    ```go
    type <ServiceName>_FooClient interface {
    	Recv() (*MsgB, error)
    	grpc.ClientStream
    }
    ```

  - **Client-Side Streaming RPC** `func YourMethod(ctx context.Context, opts ...grpc.CallOption) (YourService_YourMethodClient, error)`

    ```go
    type <ServiceName>_FooClient interface {
    	Send(*MsgA) error
    	CloseAndRecv() (*MsgB, error)
    	grpc.ClientStream
    }
    ```

  - **Bidirectional Streaming RPC** `func YourMethod(ctx context.Context, opts ...grpc.CallOption) (YourService_YourMethodClient, error)`

    ```go
    type <ServiceName>_FooClient interface {
    	Send(*MsgA) error
    	Recv() (*MsgB, error)
    	grpc.ClientStream
    }
    ```


### Try

```sh
~/ go run go_client/server.go
~/ go run go_server/client.go
```

## ALTS

gRPC authentication using Application Layer Transport Security

Note that ALTS is fully functional if the application runs on [Compute Engine](https://cloud.google.com/compute) or [Google Kubernetes Engine (GKE)](https://cloud.google.com/kubernetes-engine)

### Client

```go
import (
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/alts"
)

altsTC := alts.NewClientCreds(alts.DefaultClientOptions())
conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(altsTC))
```

### Server

```go
import (
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/alts"
)

altsTC := alts.NewServerCreds(alts.DefaultServerOptions())
server := grpc.NewServer(grpc.Creds(altsTC))
```

