package main

import (
	"context"
	"github.com/airztz/Protobuf4fun/rest/services"
	"github.com/airztz/Protobuf4fun/rest/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
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
	address  = "localhost:50051"
	rpcPort  = ":50051"
	httpPort = ":8081"
)

// server is used to implement services.HelloServer
type server struct {
	services.UnimplementedHelloServer
}

// SayHello implements services.HelloServer
func (s *server) SayHello(ctx context.Context, request *types.HelloRequest) (*types.HelloReply, error) {
	log.Printf("Received featureName: %v", request.GetFeatureName())
	featureName := request.FeatureName
	switch complexFeatureValueType := request.ComplexFeatureValue.TypeUrl; complexFeatureValueType {
	case "google.protobuf.Int32Value":
		log.Printf("Received complexValue type: %v", complexFeatureValueType)
		var x wrappers.Int32Value
		ptypes.UnmarshalAny(request.ComplexFeatureValue, &x)
		log.Printf("Received raw value: %v", x.Value)
		complexFeatureValue, _ := ptypes.MarshalAny(&wrappers.Int32Value{Value: x.Value})
		complexFeatureValue.TypeUrl = "google.protobuf.Int32Value"
		return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: complexFeatureValue}, nil
	case "google.protobuf.FloatValue":
		log.Printf("Received complexValue type: %v", complexFeatureValueType)
		var x wrappers.FloatValue
		ptypes.UnmarshalAny(request.ComplexFeatureValue, &x)
		log.Printf("Received raw value: %v", x.Value)
		complexFeatureValue, _ := ptypes.MarshalAny(&wrappers.FloatValue{Value: x.Value})
		complexFeatureValue.TypeUrl = "google.protobuf.FloatValue"
		return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: complexFeatureValue}, nil
	case "google.protobuf.BoolValue":
		log.Printf("Received complexValue type: %v", complexFeatureValueType)
		var x wrappers.BoolValue
		ptypes.UnmarshalAny(request.ComplexFeatureValue, &x)
		log.Printf("Received raw value: %v", x.Value)
		complexFeatureValue, _ := ptypes.MarshalAny(&wrappers.BoolValue{Value: x.Value})
		complexFeatureValue.TypeUrl = "google.protobuf.BoolValue"
		return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: complexFeatureValue}, nil
	case "google.protobuf.ListValue":
		log.Printf("Received complexValue type: %v", complexFeatureValueType)
		var x structpb.ListValue
		ptypes.UnmarshalAny(request.ComplexFeatureValue, &x)
		log.Printf("Received raw value: %v", x.Values)
		listValue := make([]*structpb.Value, 0)
		for _, value := range x.Values {
			switch v := value.GetKind().(type) {
			case *structpb.Value_NumberValue:
				rawValue := v.NumberValue
				log.Printf("element: %v, type: NumberValue", rawValue)
				listValue = append(listValue, &structpb.Value{Kind:&structpb.Value_NumberValue{NumberValue: rawValue}})
			case *structpb.Value_StringValue:
				rawValue := v.StringValue
				log.Printf("element: %v, type: StringValue", rawValue)
				listValue = append(listValue, &structpb.Value{Kind:&structpb.Value_StringValue{StringValue: rawValue}})
			case *structpb.Value_BoolValue:
				rawValue := v.BoolValue
				log.Printf("element: %v, type: BoolValue", rawValue)
				listValue = append(listValue, &structpb.Value{Kind:&structpb.Value_BoolValue{BoolValue: rawValue}})
			default:
				log.Printf("Failed to parse complexValue")
				return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: nil}, nil
			}
		}
		complexFeatureValue, _ := ptypes.MarshalAny(&structpb.ListValue{Values: listValue})
		complexFeatureValue.TypeUrl = "google.protobuf.ListValue"
		return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: complexFeatureValue}, nil
	case "google.protobuf.Struct":
		log.Printf("Received complexValue type: %v", complexFeatureValueType)
		var x structpb.Struct
		ptypes.UnmarshalAny(request.ComplexFeatureValue, &x)
		log.Printf("Received raw value: %v", x.Fields)
		structValue := make(map[string]*structpb.Value)
		for key, element := range x.Fields {
			switch v := element.GetKind().(type) {
			case *structpb.Value_NumberValue:
				rawValue := v.NumberValue
				log.Printf("key: %v => element: %v, type: NumberValue", key, rawValue)
				structValue[key] = &structpb.Value{Kind:&structpb.Value_NumberValue{NumberValue: rawValue}}
			case *structpb.Value_StringValue:
				rawValue := v.StringValue
				log.Printf("key: %v => element: %v, type: StringValue", key, rawValue)
				structValue[key] = &structpb.Value{Kind:&structpb.Value_StringValue{StringValue: rawValue}}
			case *structpb.Value_BoolValue:
				rawValue := v.BoolValue
				log.Printf("key: %v => element: %v, type: BoolValue", key, rawValue)
				structValue[key] = &structpb.Value{Kind:&structpb.Value_BoolValue{BoolValue: rawValue}}
			case *structpb.Value_ListValue:
				rawValue := v.ListValue
				log.Printf("key: %v => element: %v, type: ListValue", key, rawValue)
				structValue[key] = &structpb.Value{Kind:&structpb.Value_ListValue{ListValue: rawValue}}
			default:
				log.Printf("Failed to parse complexValue")
				return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: nil}, nil
			}
		}
		complexFeatureValue, _ := ptypes.MarshalAny(&structpb.Struct{Fields: structValue})
		complexFeatureValue.TypeUrl = "google.protobuf.Struct"
		return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: complexFeatureValue}, nil
	default:
		log.Printf("Failed to parse complexValue")
		return &types.HelloReply{FeatureName: featureName, ComplexFeatureValue: nil}, nil
	}
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
	err = services.RegisterHelloHandlerFromEndpoint(ctx, gatewayMux, address, []grpc.DialOption{grpc.WithInsecure()})
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
