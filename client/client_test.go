// Package client
// @Title        client_test.go
// @Description
// @Author       gxk
// @Time         2023/11/27 6:50 PM
package client

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"main/pb"
	"main/utils"
	"strconv"
	"testing"
	"time"
)

func GrpcClient(uri string) {
	//创建tcp连接
	conn, err := grpc.Dial(uri, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")), grpc.WithIdleTimeout(3*time.Second))
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return
	}
	defer conn.Close()
	//初始化user服务客户端，NewRpcGetUserInfoClient是自动生成的，格式：New{服务名}Client
	client := pb.NewUserClient(conn)
	//发起请求，注意请求方法和定义的方法一样
	resp, err := client.QueryUser(context.Background(), &pb.QueryUserRequest{UserId: 10})
	if err != nil {
		fmt.Printf("could not greet: %v\n", err)
		return
	}
	fmt.Printf("收到返回值：%+v\n", resp)
}
func Test_Grpc(T *testing.T) {

	//初始化consul配置, 客户端服务器需要一致
	consulConfig := api.DefaultConfig()
	//设置consul服务器地址: 默认127.0.0.1:8500, 如果consul部署到其它服务器上,则填写其它服务器地址
	//consulConfig.Address = "127.0.0.1:8500"
	//2、获取consul操作对象
	consulClient, _ := api.NewClient(consulConfig) //目前先屏蔽error,也可以获取error进行错误处理
	//3、获取consul服务发现地址,返回的ServiceEntry是一个结构体数组
	//参数说明:service：服务名称,服务端设置的那个Name, tag:标签,服务端设置的那个Tags,, passingOnly bool, q: 参数
	serviceEntry, queryMeta, err := consulClient.Health().Service("UserService", "grpc-gateway", true, nil)
	fmt.Printf("serviceEntry: %+v\n", utils.Json(serviceEntry))
	fmt.Printf("queryMeta: %+v\n", queryMeta)
	fmt.Printf("err: %+v\n", err)
	if len(serviceEntry) > 0 {
		//打印地址
		fmt.Println(serviceEntry[0].Service.Address)
		fmt.Println(serviceEntry[0].Service.Port)
		//拼接地址
		//strconv.Itoa: int转string型
		address := serviceEntry[0].Service.Address + ":" + strconv.Itoa(serviceEntry[0].Service.Port)
		T.Log(address)
		GrpcClient(address)
	} else {
		T.Log("未找到服务地址")
	}
}
func Test_Deregister(T *testing.T) {
	//初始化consul配置,客户端服务器需要一致
	consulConfig := api.DefaultConfig()
	//consulConfig.Address = "192.168.1.132:8500" //consul服务器的地址
	//获取consul操作对象
	registerClient, _ := api.NewClient(consulConfig)
	//注销服务ServiceDeregister(ServerID),ServerID: 微服务服务端服务发现id
	err := registerClient.Agent().ServiceDeregister("111")
	fmt.Println(err)

}
