// Package services
// @Title        user.go
// @Description
// @Author       gxk
// @Time         2023/11/27 6:22 PM
package services

import (
	"context"
	"errors"
	"main/pb"
)

type UserService struct {
	pb.UnimplementedUserServer
}

func (receiver *UserService) QueryUser(ctx context.Context, data *pb.QueryUserRequest) (*pb.QueryUserResponse, error) {
	if data.GetUserId() == 0 {
		return nil, errors.New("user_id不能为0")
	}
	return &pb.QueryUserResponse{
		UserId:   data.GetUserId(),
		UserName: "gxk",
	}, nil
}
