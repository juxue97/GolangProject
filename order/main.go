package main

import (
	"context"
	"net"
	"time"

	"github.com/juxue97/common"
	"github.com/juxue97/common/broker"
	"github.com/juxue97/common/discovery"
	"github.com/juxue97/common/discovery/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"

	jaegerAddr = common.GetString("JAEGER_ADDR", "localhost:4318")
	gRPCAddr   = common.GetString("GRPC_ADDR", "localhost:8001")
	consulAddr = common.GetString("CONSUL_ADDR", "localhost:8500")
	amqpUser   = common.GetString("RABBITMQ_USER", "juxue")
	amqpPass   = common.GetString("RABBITMQ_PASS", "veryStrongPassword")
	amqpHost   = common.GetString("RABBITMQ_HOST", "localhost")
	amqpPort   = common.GetString("RABBITMQ_PORT", "5672")
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	err = common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		logger.Fatal("failed to set global tracer", zap.Error(err))
	}

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
				logger.Error("failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	gRPCServer := grpc.NewServer()

	l, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	store := NewStore()
	service := NewService(store)
	serviceWithTelemetry := NewtelemetryMiddleware(service)
	serviceWithLogging := NewloggingMiddleware(serviceWithTelemetry)

	NewGRPCHandler(gRPCServer, serviceWithLogging, ch)

	consumer := NewConsumer(serviceWithLogging)
	go consumer.Listen(ch)

	logger.Info("gRPC server has been started at %s", zap.String("port", gRPCAddr))

	if err := gRPCServer.Serve(l); err != nil {
		logger.Fatal("failed to serve gRPC server", zap.Error(err))
	}
}
