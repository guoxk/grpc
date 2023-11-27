// Package grpc_gin
// @Title        main.go
// @Description
// @Author       gxk
// @Time         2023/11/23 5:25 PM
package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"log"
	"main/model"
	"main/pb"
	"main/services"
	"net"
	"net/http"
	"time"
)

func main() {
	go func() {
		srv := grpc.NewServer() //初始化一个grpc服务
		//RegisterRpcGetUserInfoServer这个方法是自动生成的，格式Register{服务名}Server
		pb.RegisterUserServer(srv, new(services.UserService))
		//后续可以注册进来无数个服务
		//底层通讯使用tcp
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", model.PORT))
		if err != nil {
			log.Printf("failed to listen: %v", err)
			return
		}
		if err = srv.Serve(listener); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	go func() {
		// 等待grpc服务启动
		time.Sleep(2 * time.Second)
		ctx := context.Background()
		gwMux := runtime.NewServeMux()
		// 接口超时时间
		runtime.DefaultContextTimeout = 5 * time.Second
		// 将gateway注册到grpc客户端
		err := pb.RegisterUserHandlerFromEndpoint(ctx, gwMux, fmt.Sprintf(":%d", model.PORT), []grpc.DialOption{grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip"))})
		if err != nil {
			log.Printf("RegisterUserHandlerFromEndpoint err: %v", err)
			return
		}
		// 注册http服务到gateway
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", model.HTTP_PORT),
			Handler: gwMux,
		}
		log.Printf("serving gateway on %s\n", server.Addr)
		if err = server.ListenAndServe(); err != nil {
			log.Printf("RegisterUserHandlerFromEndpoint err: %v", err)
			return
		}
	}()

	// 注册consul服务发现

	//注册consul服务
	//1、初始化consul配置
	consulConfig := api.DefaultConfig()
	//设置consul服务器地址: 默认127.0.0.1:8500, 如果consul部署到其它服务器上,则填写其它服务器地址
	//consulConfig.Address = "127.0.0.1:8500"
	//2、获取consul操作对象
	consulClient, _ := api.NewClient(consulConfig)
	// 3、配置注册服务的参数
	agentService := api.AgentServiceRegistration{
		ID:      "1",                      // 服务id,顺序填写即可
		Tags:    []string{"grpc-gateway"}, // tag标签
		Name:    "UserService",            // 服务名称, 注册到服务发现(consul)的K
		Port:    model.PORT,               // 端口号: 需要与下面的监听， 指定 IP、port一致
		Address: "127.0.0.1",              // 当前微服务部署地址: 结合Port在consul设置为V: 需要与下面的监听， 指定 IP、port一致
		Check: &api.AgentServiceCheck{ // 健康检测
			TCP:      fmt.Sprintf("127.0.0.1:%d", model.PORT), // 前微服务部署地址,端口 : 需要与下面的监听， 指定 IP、port一致
			Timeout:  "5s",                                    // 超时时间
			Interval: "3s",                                    // 循环检测间隔时间
		},
	}

	//4、注册服务到consul上
	err := consulClient.Agent().ServiceRegister(&agentService)
	if err != nil {
		log.Printf("register consul services err: %s\n", err.Error())
	} else {
		log.Printf("register consul services success\n")
	}
	select {}
}
