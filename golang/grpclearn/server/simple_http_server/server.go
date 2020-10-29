package main

import (
	pb "awesomeProject2/grpclearn/proto"
	"context"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " HTTP Server"}, nil
}

const PORT = "9003"

func main() {
	mux := GetHTTPServeMux()

	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})

	http.ListenAndServe(":"+PORT,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}

			return
		}),
	)
}

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("eddycjy: go-grpc-example"))
	})

	return mux
}
