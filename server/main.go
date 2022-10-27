package main

import (
	"log"
	"net"

	pb "github.com/jyotirmoy-paul/nsu_calculator_backend/proto"
	"github.com/jyotirmoy-paul/nsu_calculator_backend/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.CalculatorServiceServer
}

func main() {
	addr := utils.GetAddress()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	log.Printf("Listening on: %v", addr)

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &Server{})

	// register our server to work with relfection - for exploring
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
