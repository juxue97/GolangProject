package main

import (
	"context"
	"log"
	"net"

	"github.com/juxue97/common"
	"google.golang.org/grpc"
)

var gRPCAddr = common.GetString("GRPC_ADDR", "localhost:8001")

func main() {
	gRPCServer := grpc.NewServer()

	l, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	defer l.Close()

	store := NewStore()
	service := NewService(store)
	NewGRPCHandler(gRPCServer, service)

	service.createOrder(context.Background())

	log.Printf("gRPC server has been started at %s", gRPCAddr)

	if err := gRPCServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
