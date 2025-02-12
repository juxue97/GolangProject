package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/juxue97/common"
	"github.com/juxue97/common/discovery"
	"github.com/juxue97/common/discovery/consul"
	"github.com/juxue97/gateway/gateway"
)

var (
	serviceName = "gateway"

	jaegerAddr = common.GetString("JAEGER_ADDR", "localhost:4318")
	httpAddr   = common.GetString("HTTP_ADDR", ":8080")
	consulAddr = common.GetString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		log.Fatal("failed to set global tracer")
	}

	// register service on consul
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
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

	// expose http server here, then grpc to other services
	ordersGateway := gateway.NewOrdersGateway(registry)
	stocksGateway := gateway.NewStocksGateway(registry)

	mux := http.NewServeMux()
	handler := NewHandler(ordersGateway, stocksGateway)
	handler.registerRoutes(mux)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start the http server")
	}
}
