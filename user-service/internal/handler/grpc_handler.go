package handler

import (
	"context"
	"fmt"
	"user-service/internal/entity"
	"user-service/internal/middleware"
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

func (h *grpcHandler) GetProfile(ctx context.Context, _ *userGRPC.GetProfileRequest) (*userGRPC.GetProfileResponse, error) {
	// Извлекаем userID из контекста (установлен middleware'ом)
	userID, err := middleware.ExtractUserID(ctx)
	if err != nil {
		return nil, err
	}

	userInfo, err := h.userService.GetProfileInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &userGRPC.GetProfileResponse{
		User: &userGRPC.User{
			Id:          userInfo.ID,
			FirstName:   userInfo.FirstName,
			LastName:    userInfo.LastName,
			Email:       userInfo.Email,
			PhoneNumber: userInfo.PhoneNumber,
			CreatedAt:   timestamppb.New(userInfo.CreatedAt),
			UpdatedAt:   timestamppb.New(userInfo.UpdatedAt),
		},
	}, nil
}

func (h *grpcHandler) Logout(ctx context.Context, _ *userGRPC.LogoutRequest) (*userGRPC.LogoutResponse, error) {
	userId, err := middleware.ExtractUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = h.userService.Logout(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &userGRPC.LogoutResponse{Success: true}, nil
}

func (h *grpcHandler) UpdateProfile(ctx context.Context, req *userGRPC.UpdateProfileRequest) (*userGRPC.UpdateProfileResponse, error) {
	userID, err := middleware.ExtractUserID(ctx)
	if err != nil {
		return nil, err
	}
	updateUser, err := h.userService.UpdateUserProfile(ctx, &entity.User{
		ID:          userID,
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
	})
	if err != nil {
		return nil, err
	}

	return &userGRPC.UpdateProfileResponse{
		User: &userGRPC.User{
			Id:          updateUser.ID,
			FirstName:   updateUser.FirstName,
			LastName:    updateUser.LastName,
			Email:       updateUser.Email,
			PhoneNumber: updateUser.PhoneNumber,
			CreatedAt:   timestamppb.New(updateUser.CreatedAt),
			UpdatedAt:   timestamppb.New(updateUser.UpdatedAt),
		},
	}, nil

}
