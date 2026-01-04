package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func JWTInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("JWT Interceptor: method=%s, secret_len=%d", info.FullMethod, len(secret))

		// Пропускаем reflection запросы без аутентификации
		if strings.HasPrefix(info.FullMethod, "/grpc.reflection.") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("JWT Interceptor: no metadata")
			return nil, status.Errorf(codes.Unauthenticated, "no metadata")
		}

		log.Printf("JWT Interceptor: metadata keys=%v", md)

		auth := md.Get("authorization")
		if len(auth) == 0 {
			log.Printf("JWT Interceptor: no token, auth=%v", auth)
			return nil, status.Errorf(codes.Unauthenticated, "no token")
		}

		// Поддержка токена как с "Bearer " префиксом так и без
		tokenStr := auth[0]
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		}
		log.Printf("JWT Interceptor: token=%s...", tokenStr[:min(20, len(tokenStr))])

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			log.Printf("JWT Interceptor: invalid token, err=%v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		log.Printf("JWT Interceptor: claims=%v", claims)

		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Println("JWT Interceptor: no user ID in claims")
			return nil, status.Errorf(codes.Unauthenticated, "no user ID")
		}

		log.Printf("JWT Interceptor: userID=%s", userID)
		ctx = context.WithValue(ctx, "userID", userID)
		return handler(ctx, req)
	}
}
