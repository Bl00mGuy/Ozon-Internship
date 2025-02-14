package server

import (
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/Bl00mGuy/url-shortener/internal/repository"
	"github.com/Bl00mGuy/url-shortener/internal/services"
	pb "github.com/Bl00mGuy/url-shortener/proto/gen/go"
)

func InitializeGRPCServer(port string, repo repository.URLStorage, logger *logrus.Logger) {
	address := ":" + port

	lis, err := configureGRPCListener(address)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	err = registerGRPCService(grpcServer, repo, logger)
	if err != nil {
		logger.Fatalf("failed to register gRPC service: %v", err)
		return
	}

	err = runGRPCServer(grpcServer, lis)
	if err != nil {
		logger.Fatalf("failed to serve gRPC: %v", err)
	}
}

func configureGRPCListener(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}

func registerGRPCService(grpcServer *grpc.Server, repo repository.URLStorage, logger *logrus.Logger) error {
	pb.RegisterUrlManagerServer(grpcServer, services.NewUrlShortener(repo, logger))
	return nil
}

func runGRPCServer(grpcServer *grpc.Server, lis net.Listener) error {
	return grpcServer.Serve(lis)
}
