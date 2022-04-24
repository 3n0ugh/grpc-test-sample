package main

import (
	"log"
	"net"

	"github.com/3n0ugh/grpc-test-sample/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	baseServer := grpc.NewServer()
	pb.RegisterTelephoneServer(baseServer, NewServer())

	log.Printf("Starting server on %s", lis.Addr().String())
	baseServer.Serve(lis)
}
