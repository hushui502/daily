package main

import (
	pb "awesomeProject2/grpclearn/proto"
	"awesomeProject2/procase"
	"io"
	"log"
)

type StreamService struct{}

func (s StreamService) List(*procase.StreamRequest, procase.StreamService_ListServer) error {
	panic("implement me")
}

func (s StreamService) Record(procase.StreamService_RecordServer) error {
	panic("implement me")
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gPRC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}

const PORT = "9001"
