package main

import (
	"log"
	"net"

	"github.com/upamune/grpc-error-sample/proto"
	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	api.RegisterUserServiceServer(s, NewUserService())

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server is running...")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
