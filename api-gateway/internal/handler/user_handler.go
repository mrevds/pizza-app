package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrevds/pizza-app/api-gateway/internal/client"
	"go.uber.org/zap"
)

type UserHandler struct {
	clients *client.GRPCClients
	logger  *zap.Logger
}

func NewUserHandler(clients *client.GRPCClients, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		clients: clients,
		logger:  logger,
	}
}

// Health check endpoint
func (h *UserHandler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "api-gateway"})
}

// ...existing code...

// GetProfile - получить профиль пользователя
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем токен из контекста (установлен в middleware)
	token := r.Context().Value("token").(string)

	h.logger.Info("fetching user profile",
		zap.String("token", token[:20]+"..."),
	)

	// Здесь вы будете вызывать gRPC User Service
	// Пока просто заглушка
	response := map[string]interface{}{
		"id":    "123",
		"name":  "John Doe",
		"phone": "+1234567890",
	}

	_ = json.NewEncoder(w).Encode(response)
}

// Register - регистрация пользователя
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		FirstName   string `json:"first_name"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("user registration",
		zap.String("phone", req.PhoneNumber),
	)

	// Здесь вы будете вызывать gRPC User Service
	response := map[string]interface{}{
		"id":    "456",
		"name":  req.FirstName,
		"phone": req.PhoneNumber,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
