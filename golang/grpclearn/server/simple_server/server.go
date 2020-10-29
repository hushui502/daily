package main

import (
	"awesomeProject2/grpclearn/interceptor"
	pb "awesomeProject2/grpclearn/proto"
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"log"
	"net"
)

type SearchService struct {
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	opts := []grpc.ServerOption{
		//grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			interceptor.RecoverInterceptor,
			interceptor.LoggingInterceptor,
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}
