package main

import (
	"fmt"
	"log"
	"net"

	"bibit.id/challenge/handler"
	"bibit.id/challenge/proto"
	"google.golang.org/grpc"
)

func serveGRPC(grpcHandler *handler.Handler) {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterBibitServer(grpcServer, grpcHandler)

	fmt.Println("GRPC server: listening on port 50051")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
