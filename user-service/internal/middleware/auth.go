package middleware

import (
	"context"
	"user-service/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager *utils.JWTManager
}

func NewAuthInterceptor(jwtManager *utils.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Список методов, которые НЕ требуют авторизации
		publicMethods := map[string]bool{
			"/user_service.v1.UserService/Register":      true,
			"/user_service.v1.UserService/Login":         true,
			"/user_service.v1.UserService/RefreshTokens": true,
		}

		// Если метод публичный - пропускаем без проверки
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Извлекаем metadata из контекста
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		// Получаем access token из заголовка authorization
		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		accessToken := values[0]

		// Убираем префикс "Bearer" если есть
		if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
			accessToken = accessToken[7:]
		}

		// Validate access token
		claims, err := a.jwtManager.ValidateToken(accessToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid access token: %v", err)
		}

		// Проверяем что это access token, а не refresh
		if claims.Type != "access" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token type")
		}

		// Добавляем user_id в контекст для использования в handlers
		ctx = context.WithValue(ctx, "user_id", claims.UserID)

		// call the handler
		return handler(ctx, req)
	}
}

// ExtractUserID - безопасно извлекает userID из контекста
func ExtractUserID(ctx context.Context) (string, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return "", status.Errorf(codes.Unauthenticated, "user_id not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", status.Errorf(codes.Internal, "invalid user_id type")
	}

	return userIDStr, nil
}
