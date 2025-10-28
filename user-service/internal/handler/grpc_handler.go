package handler

import (
	"context"
	"fmt"
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

func (h *grpcHandler) Login(ctx context.Context, req *userGRPC.LoginRequest) (*userGRPC.LoginResponse, error) {
	user, accessToken, refreshToken := h.userService.Login(ctx, req.PhoneNumber, req.Password)
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &userGRPC.LoginResponse{
		User: &userGRPC.User{
			Id:          user.ID,
			FirstName:   user.FirstName,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *grpcHandler) RefreshTokens(ctx context.Context, req *userGRPC.RefreshTokensRequest) (*userGRPC.RefreshTokensResponse, error) {
	accessToken, refreshToken, err := h.userService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &userGRPC.RefreshTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
