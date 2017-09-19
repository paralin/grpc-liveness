package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/paralin/grpc-liveness/statussvc"
	"google.golang.org/grpc"
)

const (
	port = ":5000"
)

// server is used to implement statussvc.StatusService.
type server struct{}

var isReady bool

// GetReadiness checks if the service is ready. Returns an error if not ready.
func (s *server) GetReadiness(
	ctx context.Context,
	in *statussvc.GetReadinessRequest,
) (*statussvc.GetReadinessResponse, error) {
	if !isReady {
		return nil, errors.New("not ready yet")
	}
	return &statussvc.GetReadinessResponse{}, nil
}

// GetLiveness checks if the service is alive at all. Returns an error if not alive.
func (s *server) GetLiveness(
	ctx context.Context,
	req *statussvc.GetLivenessRequest,
) (*statussvc.GetLivenessResponse, error) {
	return &statussvc.GetLivenessResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	statussvc.RegisterStatusServiceServer(s, &server{})

	delayTime := 5
	log.Printf("becoming ready in %d seconds!", delayTime)
	go func() {
		time.Sleep(time.Duration(delayTime) * time.Second)
		log.Printf("becoming ready!")
		isReady = true
	}()
	log.Printf("listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
