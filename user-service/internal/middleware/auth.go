package middleware

import (
	"context"
	"user-service/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor - middleware для проверки access token
type AuthInterceptor struct {
	jwtManager *utils.JWTManager
}

func NewAuthInterceptor(jwtManager *utils.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}

// Unary - перехватчик для unary gRPC методов
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

		// Валидируем access token
		claims, err := a.jwtManager.ValidateToken(accessToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid access token: %v", err)
		}

		// Проверяем что это access token, а не refresh
		if claims.Type != "access" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token type")
		}

		// Добавляем user_id в контекст для использования в хендлерах
		ctx = context.WithValue(ctx, "user_id", claims.UserID)

		// Вызываем хендлер
		return handler(ctx, req)
	}
}
