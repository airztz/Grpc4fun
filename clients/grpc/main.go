package main

import (
	"context"
	"github.com/airztz/Protobuf4fun/grpc/services"
	"github.com/airztz/Protobuf4fun/grpc/types"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := services.NewHelloClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	t1 := &timestamp.Timestamp{
		Seconds: 5, // easy to verify
		Nanos:   6, // easy to verify
	}
	//t2 := &
	serialized, err := proto.Marshal(t1)
	if err != nil {
		log.Fatal("could not serialize timestamp")
	}
	r, err := client.SayHello(ctx, &types.HelloRequest{ComplexFeatureValue: &any.Any{
		TypeUrl: proto.MessageName(t1),
		Value:   serialized,
	}})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.ComplexFeatureValue)
}
