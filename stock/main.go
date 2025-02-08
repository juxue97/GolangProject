package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload" // put this line for all modules
	"github.com/juxue97/common"
	"github.com/juxue97/common/broker"
	mongoConn "github.com/juxue97/common/db"
	"github.com/juxue97/common/discovery"
	"github.com/juxue97/common/discovery/consul"
	"github.com/juxue97/stock/gateway"
	stripeProcessor "github.com/juxue97/stock/processor/stripe"
	"github.com/stripe/stripe-go/v81"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "stocks"

	gRPCAddr = common.GetString("GRPC_ADDR", "localhost:8003")
	httpAddr = common.GetString("HTTP_ADDR", ":8083")

	jaegerAddr = common.GetString("JAEGER_ADDR", "localhost:4318")
	consulAddr = common.GetString("CONSUL_ADDR", "localhost:8500")

	amqpUser = common.GetString("RABBITMQ_USER", "juxue")
	amqpPass = common.GetString("RABBITMQ_PASS", "veryStrongPassword")
	amqpHost = common.GetString("RABBITMQ_HOST", "localhost")
	amqpPort = common.GetString("RABBITMQ_PORT", "5672")

	mongoUser = common.GetString("MONGO_DB_USER", "juxue")
	mongoPass = common.GetString("MONGO_DB_PASS", "veryStrongPassword")
	mongoHost = common.GetString("MONGO_DB_HOST", "localhost:27017")

	stripeKey = common.GetString("STRIPE_KEY", "")
	// endpointStripeSecret = common.GetString("ENDPOINT_STRIPE_SECRET", "whsec_...")
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

	stripe.Key = stripeKey

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPass, mongoHost)
	mongoClient, err := mongoConn.ConnectToMongoDB(mongoURI)
	if err != nil {
		logger.Fatal("failed to connect to mongo", zap.Error(err))
	}
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	gRPCServer := grpc.NewServer()

	l, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	stripeProcessor := stripeProcessor.NewProcessor()
	gateway := gateway.NewGateway(registry)

	store := NewStore(mongoClient)
	service := NewStockService(store, stripeProcessor, gateway)
	serviceWithTelemetry := NewTelemetryMiddleware(service)
	serviceWithLogging := NewLoggingMiddleware(serviceWithTelemetry)

	mux := http.NewServeMux()
	httpServer := NewStockHandler(service)
	httpServer.registerRouters(mux)

	// adding prefix
	v1 := http.NewServeMux()
	v1.Handle("/api/", http.StripPrefix("/api", mux))

	go func() {
		logger.Info("starting http server", zap.String("port", httpAddr))
		if err := http.ListenAndServe(httpAddr, v1); err != nil {
			logger.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	NewGRPCHandler(gRPCServer, serviceWithLogging, ch)

	consumer := NewConsumer(serviceWithLogging)
	go consumer.Listen(ch)

	logger.Info("gRPC server has been started at %s", zap.String("port", gRPCAddr))

	if err := gRPCServer.Serve(l); err != nil {
		logger.Fatal("failed to serve gRPC server", zap.Error(err))
	}
}
