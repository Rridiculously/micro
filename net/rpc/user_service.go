package rpc

import (
	"context"
	"log"
)

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "UserService"
}

type GetByIdReq struct {
	Id int
}
type GetByIdResp struct {
	Msg string
}

type UserServiceServer struct {
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	log.Println(req)
	return &GetByIdResp{
		Msg: "hello,world",
	}, nil
}
func (u *UserServiceServer) Name() string {
	return "UserService"
}
