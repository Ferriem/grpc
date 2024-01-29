package main

import (
	"context"
	"flag"
	"fmt"
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
