package main

import (
	"context"
	"net"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload" // put this line for all modules
	"github.com/juxue97/common"
	"github.com/juxue97/common/broker"
	"github.com/juxue97/common/discovery"
	"github.com/juxue97/common/discovery/consul"
	"github.com/juxue97/payment/gateway"
	stripeProcessor "github.com/juxue97/payment/processor/stripe"
	"github.com/stripe/stripe-go/v81"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "payments"

	httpAddr = common.GetString("HTTP_ADDR", ":8079")

	jaegerAddr           = common.GetString("JAEGER_ADDR", "localhost:4318")
	gRPCAddr             = common.GetString("GRPC_ADDR", "localhost:8002")
	consulAddr           = common.GetString("CONSUL_ADDR", "localhost:8500")
	amqpUser             = common.GetString("RABBITMQ_USER", "juxue")
	amqpPass             = common.GetString("RABBITMQ_PASS", "veryStrongPassword")
	amqpHost             = common.GetString("RABBITMQ_HOST", "localhost")
	amqpPort             = common.GetString("RABBITMQ_PORT", "5672")
	stripeKey            = common.GetString("STRIPE_KEY", "")
	endpointStripeSecret = common.GetString("ENDPOINT_STRIPE_SECRET", "whsec_...")
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

	// stripe conn
	stripe.Key = stripeKey

	// HTTPServer
	mux := http.NewServeMux()
	httpServer := NewPaymentHTTPHandler(ch)
	httpServer.registerRouters(mux)

	go func() {
		logger.Info("starting http server", zap.String("port", httpAddr))
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			logger.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	// grpcServer
	gRPCServer := grpc.NewServer()

	l, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	stripeProcessor := stripeProcessor.NewProcessor()
	gateway := gateway.NewGateway(registry)

	service := NewPaymentService(stripeProcessor, gateway)
	serviceWithTelemetry := NewtelemetryMiddleware(service)
	serviceWithLogging := NewloggingMiddleware(serviceWithTelemetry)

	amqpConsumer := NewConsumer(serviceWithLogging)
	// listen method
	go amqpConsumer.Listen(ch)

	// store := NewStore()
	// service := NewService(store)
	// NewGRPCHandler(gRPCServer, service, ch)

	// service.createOrder(context.Background())

	logger.Info("gRPC server has been started", zap.String("port", gRPCAddr))

	if err := gRPCServer.Serve(l); err != nil {
		logger.Fatal("failed to serve gRPC server", zap.Error(err))
	}
}
