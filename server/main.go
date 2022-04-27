package main

import (
	"context"
	"fmt"
	pb "github.com/fukaraca/gRPC-hello-Mars/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedMessagingServiceServer
}

func (s *server) SendMessage(ctx context.Context, request *pb.CreateMessageRequest) (*pb.CreateMessageResponse, error) {
	if request.Request.Text != "" {
		temp := &pb.MessageResponse{Read: true}
		fmt.Printf("the message: %s \nfrom: %s \nreceived at %s\n\n", request.Request.Text, request.Request.Sender, request.Request.SendingTime.AsTime().String())
		return &pb.CreateMessageResponse{Response: temp}, nil
	}
	return nil, fmt.Errorf("no text has been received")
}

func (s *server) MustEmbedUnimplementedMessagingServiceServer() {

}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("listening on port 8080 failed", err)
	}
	s := grpc.NewServer()
	pb.RegisterMessagingServiceServer(s, &server{})

	err = s.Serve(listener)
	if err != nil {
		log.Fatalln("listening failed on serve", err)
	}
}
