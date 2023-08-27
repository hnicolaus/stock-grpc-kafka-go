package main

import (
	"log"
	"net"

	"bibit.id/challenge/handler"
	"bibit.id/challenge/proto"
	"google.golang.org/grpc"
)

func serveGRPC(grpcHandler *handler.Handler) {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Printf("[GRPC] Failed to listen to port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterBibitServer(grpcServer, grpcHandler)

	log.Print("[GRPC] Listening on port 50051")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("[GRPC] Failed to serve GRPC server: %v", err)
	}
}
