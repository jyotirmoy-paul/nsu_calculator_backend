package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/jyotirmoy-paul/nsu_calculator_backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Operate(ctx context.Context, in *pb.OperationRequest) (*pb.OperationResponse, error) {
	var result float64

	switch in.OperationType {
	case pb.OperationType_OPERATION_TYPE_ADD:
		result = in.OperandA + in.OperandB
	case pb.OperationType_OPERATION_TYPE_SUBTRACT:
		result = in.OperandA - in.OperandB
	case pb.OperationType_OPERATION_TYPE_MULTIPLY:
		result = in.OperandA * in.OperandB
	case pb.OperationType_OPERATION_TYPE_DIVIDE:
		if in.OperandB == 0 {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"Cannot divide by Zero",
			)
		}
		result = in.OperandA / in.OperandB
	default:
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid Operation Type: %v", in.OperationType),
		)
	}

	return &pb.OperationResponse{
		Result: result,
	}, nil
}

func (s *Server) Factorize(in *pb.FactorizationRequest, stream pb.CalculatorService_FactorizeServer) error {
	var k int32 = 2
	N := in.Number

	for N > 1 {
		if N%k == 0 {
			stream.Send(&pb.FactorizationResponse{
				Factor: k,
			})
			N = N / k
		} else {
			k = k + 1
		}
	}

	return nil
}

func (s *Server) FindAverage(stream pb.CalculatorService_FindAverageServer) error {
	var sum float64
	var count int

	for {
		msg, err := stream.Recv()

		// Client has stopped streaming, let's return the result
		if err == io.EOF {

			var result float64

			if count == 0 {
				result = -1
			} else {
				result = sum / float64(count)
			}

			return stream.SendAndClose(&pb.AverageResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Printf("Error in FindAverage stream: %v\n", err)
			return status.Errorf(
				codes.Unknown,
				"Something went wrong",
			)
		}

		sum += msg.Number
		count++
	}
}

func (s *Server) Sum(stream pb.CalculatorService_SumServer) error {
	var sum float64

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Printf("Error in FindMax stream: %v\n", err)
			return status.Errorf(
				codes.Unknown,
				"Something went wrong",
			)
		}

		sum += msg.Number

		stream.Send(&pb.SumResponse{
			Number: sum,
		})
	}
}
