package handler

import (
	"user-service/internal/service"
	userGRPC "user-service/pkg/user-service_v1"
)

type grpcHandler struct {
	userGRPC.UnimplementedUserServiceServer
	userService service.UserService
}

func NewGRPCHandler(s service.UserService) userGRPC.UserServiceServer {
	return &grpcHandler{userService: s}
}
