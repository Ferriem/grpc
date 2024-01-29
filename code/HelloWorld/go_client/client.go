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

}
