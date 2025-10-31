package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrevds/pizza-app/api-gateway/internal/client"
	"github.com/mrevds/pizza-app/api-gateway/internal/handler"
	"github.com/mrevds/pizza-app/api-gateway/internal/middleware"
	"go.uber.org/zap"
)

func NewRouter(clients *client.GRPCClients, logger *zap.Logger) *mux.Router {
	router := mux.NewRouter()

	// Global middleware
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))

	// Health check
	userHandler := handler.NewUserHandler(clients, logger)
	router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")

	// ===== AUTH ROUTES (без аутентификации) =====
	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", userHandler.Register).Methods("POST")   // Пока заглушка
	authRouter.HandleFunc("/refresh", userHandler.Register).Methods("POST") // Пока заглушка

	// ===== USER ROUTES (требует аутентификацию) =====
	userProtectedRouter := router.PathPrefix("/api/v1/user").Subrouter()
	userProtectedRouter.Use(middleware.AuthMiddleware)
	userProtectedRouter.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")
	userProtectedRouter.HandleFunc("/profile", userHandler.Register).Methods("PUT") // Пока заглушка
	userProtectedRouter.HandleFunc("/logout", userHandler.Register).Methods("POST") // Пока заглушка

	// ===== CARD ROUTES (требует аутентификацию) =====
	cardHandler := handler.NewCardHandler(clients, logger)
	cardRouter := router.PathPrefix("/api/v1/cards").Subrouter()
	cardRouter.Use(middleware.AuthMiddleware)
	cardRouter.HandleFunc("", cardHandler.GetUserCards).Methods("GET")
	cardRouter.HandleFunc("", cardHandler.AddCard).Methods("POST")
	cardRouter.HandleFunc("/balance", cardHandler.GetBalance).Methods("GET")
	cardRouter.HandleFunc("/deposit", cardHandler.Deposit).Methods("POST")
	cardRouter.HandleFunc("/withdraw", cardHandler.Deposit).Methods("POST") // Пока заглушка
	cardRouter.HandleFunc("/transfer", cardHandler.Deposit).Methods("POST") // Пока заглушка

	// 404 handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"route not found"}`))
	})

	return router
}
