// Package grpc_gin
// @Title        main.go
// @Description
// @Author       gxk
// @Time         2023/11/23 5:25 PM
package main

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"log"
	"main/pb"
	"net"
	"net/http"
	"time"
)

type UserService struct {
	pb.UnimplementedUserServer
}

func (receiver *UserService) QueryUser(ctx context.Context, data *pb.QueryUserRequest) (*pb.QueryUserResponse, error) {
	if data.GetUserId() == 0 {
		return nil, errors.New("user_id不能为0")
	}
	//time.Sleep(5 * time.Second)
	return &pb.QueryUserResponse{
		UserId:   data.GetUserId(),
		UserName: "gxk",
	}, nil
}

const PORT = ":1234"
const HTTP_PORT = ":8080"

func main() {
	go func() {
		srv := grpc.NewServer() //初始化一个grpc服务
		//RegisterRpcGetUserInfoServer这个方法是自动生成的，格式Register{服务名}Server
		pb.RegisterUserServer(srv, new(UserService))
		//后续可以注册进来无数个服务
		//底层通讯使用tcp
		listener, err := net.Listen("tcp", PORT)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err = srv.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// 等待服务启动
	time.Sleep(2 * time.Second)
	//创建tcp连接
	ctx := context.Background()
	gwMux := runtime.NewServeMux()
	// 接口超时时间
	runtime.DefaultContextTimeout = 5 * time.Second
	// 将gateway注册到grpc客户端
	err := pb.RegisterUserHandlerFromEndpoint(ctx, gwMux, PORT, []grpc.DialOption{grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip"))})
	if err != nil {
		log.Fatalf("RegisterUserHandlerFromEndpoint err: %v", err)
	}
	// 注册http服务到gateway
	server := http.Server{
		Addr:    HTTP_PORT,
		Handler: gwMux,
	}
	log.Printf("serving gateway on %s\n", server.Addr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("RegisterUserHandlerFromEndpoint err: %v", err)
	}
}
