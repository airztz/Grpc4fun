package main

import (
	"context"
	"github.com/airztz/Protobuf4fun/rest/services"
	"github.com/airztz/Protobuf4fun/rest/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	address = "localhost:50051"
	rpcPort = ":50051"
	httpPort = ":8081"
)

// server is used to implement services.HelloServer
type server struct {
	services.UnimplementedHelloServer
}

// SayHello implements services.HelloServer
func (s *server) SayHello(ctx context.Context, request *types.HelloRequest) (*types.HelloReply, error) {
	log.Printf("Received: %v", request.GetName())
	return &types.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	rpcServer := grpc.NewServer()
	services.RegisterHelloServer(rpcServer, &server{})

	gatewayMux := runtime.NewServeMux()
	err = services.RegisterHelloHandlerFromEndpoint(ctx, gatewayMux,  address, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("failed to register rpc server to gateway server: %v", err)
	}
	httpServer := &http.Server{Addr: httpPort, Handler: gatewayMux}


	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			log.Fatalf("rpc server failed to serve: %v", err)
		}
	}()

	go func() {
		// Start HTTP server (and proxy calls to gRPC server endpoint)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("http server failed to serve: %v", err)
		}
	}()

	<-shutdown
	rpcServer.GracefulStop()
	err = httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failed to stop http server: %v", err)
	}

}
