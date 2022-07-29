package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-learn/pb"
	"log"
	"net"
	"net/http"
)

var (
	tcpPort  = flag.Int("tcpPort", 8080, "The tcp server Port")
	httpPort = flag.Int("httpPort", 8081, "The http server Port")
)

// server is used to implement pb.WaylonServer.
type server struct {
	pb.UnimplementedWaylonServer
}

func (s server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println(fmt.Sprintf("Received: %s", request.GetName()))
	return &pb.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func main() {
	flag.Parse()
	// Create a listener on TCP Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *tcpPort))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// 创建一个gRPC server对象
	s := grpc.NewServer()
	// 注册RegisterWaylonServer
	pb.RegisterWaylonServer(s, &server{})
	// port端口启动gRPC Server
	log.Println("Serving gRPC on 0.0.0.0" + fmt.Sprintf(":%d", *tcpPort))
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// 创建一个连接到我们刚刚启动的 gRPC 服务器的客户端连接
	// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", *tcpPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// 注册Greeter
	err = pb.RegisterWaylonHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: gwmux,
	}
	// httpPort端口提供gRPC-Gateway服务
	log.Println("Serving gRPC-Gateway on http://0.0.0.0" + fmt.Sprintf(":%d", *httpPort))
	log.Fatalln(gwServer.ListenAndServe())
}
