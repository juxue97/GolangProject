package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/juxue97/common"
	"github.com/juxue97/common/discovery"
	"github.com/juxue97/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"

	gRPCAddr   = common.GetString("GRPC_ADDR", "localhost:8001")
	consulAddr = common.GetString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, gRPCAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

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
