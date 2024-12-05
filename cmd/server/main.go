package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"train-booking-service/pkg/adapter"
	"train-booking-service/pkg/service"
	"train-booking-service/proto"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

func main() {

	// Initialize the DAO and the service
	dao := adapter.NewTrainDAO()
	trainService := service.NewTrainService(dao)

	// Set up the gRPC server
	server := grpc.NewServer()
	proto.RegisterTrainServiceServer(server, trainService)
	reflection.Register(server)

	listener, err := net.Listen("tcp", ":7001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Listen for termination signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down gracefully...")
	server.GracefulStop()
}
