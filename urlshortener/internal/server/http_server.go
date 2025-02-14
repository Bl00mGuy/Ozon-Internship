package server

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Bl00mGuy/url-shortener/proto/gen/go"
)

func InitializeHTTPGateway(grpcPort, httpPort string, logger *logrus.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux, err := configureHTTPMux(grpcPort, ctx, logger)
	if err != nil {
		logger.Fatalf("failed to create HTTP mux: %v", err)
		return
	}

	err = runHTTPServer(mux, httpPort, logger)
	if err != nil {
		logger.Fatalf("failed to start HTTP server: %v", err)
	}
}

func configureHTTPMux(grpcPort string, ctx context.Context, logger *logrus.Logger) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterUrlManagerHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		logger.Errorf("failed to register HTTP handler: %v", err)
		return nil, err
	}
	return mux, nil
}

func runHTTPServer(mux *runtime.ServeMux, httpPort string, logger *logrus.Logger) error {
	logger.Infof("HTTP server is running on port %s", httpPort)
	return http.ListenAndServe(":"+httpPort, mux)
}
