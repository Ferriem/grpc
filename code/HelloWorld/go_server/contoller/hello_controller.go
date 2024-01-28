package contoller

import (
	"fmt"

	"github.com/Ferriem/grpc/code/HelloWorld/proto/hello"

	"golang.org/x/net/context"
)

type HelloController struct{}

func (h *HelloController) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	fmt.Println("SayHello")
	return &hello.HelloResponse{Message: fmt.Sprintf("%s", in.Name)}, nil
}

func (h *HelloController) LotsOfReplies(in *hello.HelloRequest, stream hello.Hello_LotsOfRepliesServer) error {
	fmt.Println("LotsOfReplies")
	for i := 0; i < 10; i++ {
		stream.Send(&hello.HelloResponse{Message: fmt.Sprintf("%s %s %d", in.Name, "Reply", i)})
	}
	return nil
}
