package handler

import (
	"context"
	"user-service/internal/service"
	userGRPC "user-service/pkg/user-service_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcHandler struct {
	userGRPC.UnimplementedUserServiceServer
	userService service.UserService
}

func NewGRPCHandler(s service.UserService) userGRPC.UserServiceServer {
	return &grpcHandler{userService: s}
}

func (h *grpcHandler) Register(ctx context.Context, req *userGRPC.RegisterRequest) (*userGRPC.RegisterResponse, error) {
	user, err := h.userService.Register(ctx, service.RegisterInput{
		FirstName:   req.FirstName,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &userGRPC.RegisterResponse{
		User: &userGRPC.User{
			Id:          user.ID,
			FirstName:   user.FirstName,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
