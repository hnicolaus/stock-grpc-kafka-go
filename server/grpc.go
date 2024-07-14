/*
	Hans Nicolaus
	29 Aug 2023
*/

package server

import (
	"log"
	"net"

	"stock/handler"
	"stock/model"
	"stock/proto"

	"google.golang.org/grpc"
)

func ServeGRPC(cfg model.Config, grpcHandler *handler.Handler) {
	listen, err := net.Listen(cfg.GRPC.Network, cfg.GRPC.Port)
	if err != nil {
		log.Printf("[GRPC] Failed to listen to port %s: %v", cfg.GRPC.Port, err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStockServer(grpcServer, grpcHandler)

	log.Printf("[GRPC] Serving on port %v", cfg.GRPC.Port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("[GRPC] Failed to serve GRPC server: %v", err)
	}
}
