package main

import (
	"awesomeProject2/grpclearn/interceptor"
	pb "awesomeProject2/grpclearn/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type StreamService struct{}

const (
	PORT = "9002"
)

func main() {
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			interceptor.RecoverInterceptor,
			interceptor.LoggingInterceptor,
		),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n < 4; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %v, pt.value: %v", r.Pt.Name, r.Pt.Value)
	}
	return nil
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	return nil
}
