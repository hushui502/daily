package main

import (
	fmt "fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var (
	port = ":5000"
)

type myGrpcServer struct{}

func (s *myGrpcServer) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing credentials")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["login"]; ok {
		appid = val[0]
	}
	if val, ok := md["password"]; ok {
		appkey = val[0]
	}

	if appid != "gopher" || appkey != "password" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token: appid=%s, appkey=%s", appid, appkey)
	}

	return &HelloReply{Message: "Hello " + in.Name}, nil
}

type Authentication struct {
	Login    string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"login": a.Login, "password": a.Password}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func main() {
	go startServer()
	time.Sleep(time.Second)

	doClientWork()
}

func startServer() {
	server := grpc.NewServer(grpc.UnaryInterceptor(filter))
	RegisterGreeterServer(server, new(myGrpcServer))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("could not list on %s: %s", port, err)
	}

	if err := server.Serve(lis); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}
}

func startServer1() {
	creds, err := credentials.NewServerTLSFromFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	RegisterGreeterServer(grpcServer, new(myGrpcServer))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	})
	http.ListenAndServeTLS(port, "", "", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}))

}

func filter(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
	log.Println("filter: ", info)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic : %v", r)
		}
	}()

	return handler(ctx, req)
}

func doClientWork() {
	auth := Authentication{
		Login:    "gopher",
		Password: "password",
	}

	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &HelloRequest{Name: "gopher"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("doClientWork: %s", r.Message)
}
