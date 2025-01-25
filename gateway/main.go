package main

import (
	"log"
	"net/http"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr         = common.GetString("HTTP_ADDR", ":8080")
	orderServiceAddr = "localhost:8001"
)

func main() {
	// expose http server here, then grpc to other services
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		log.Fatalf("fail to dial server: %v", err)
	}
	defer conn.Close()

	log.Println("gRPC server has been started at ", orderServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start the http server")
	}
}
