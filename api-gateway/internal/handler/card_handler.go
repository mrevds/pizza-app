package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrevds/pizza-app/api-gateway/internal/client"
	"go.uber.org/zap"
)

type CardHandler struct {
	clients *client.GRPCClients
	logger  *zap.Logger
}

func NewCardHandler(clients *client.GRPCClients, logger *zap.Logger) *CardHandler {
	return &CardHandler{
		clients: clients,
		logger:  logger,
	}
}

// GetUserCards - получить все карты пользователя
func (h *CardHandler) GetUserCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	h.logger.Info("fetching user cards", zap.String("user_id", userID))

	// Здесь вы будете вызывать gRPC Card Service
	response := map[string]interface{}{
		"cards": []map[string]interface{}{},
	}

	_ = json.NewEncoder(w).Encode(response)
}

// AddCard - добавить новую карту
func (h *CardHandler) AddCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		UserID         int64  `json:"user_id"`
		CardNumber     string `json:"card_number"`
		CardHolderName string `json:"card_holder_name"`
		ExpiryDate     string `json:"expiry_date"`
		CVV            string `json:"cvv"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("adding card",
		zap.Int64("user_id", req.UserID),
		zap.String("card_number", req.CardNumber[:4]+"..."),
	)

	// Здесь вы будете вызывать gRPC Card Service
	response := map[string]interface{}{
		"id":                 "card-123",
		"user_id":            req.UserID,
		"card_number_masked": "4532 **** **** 9010",
		"card_holder_name":   req.CardHolderName,
		"balance":            0.0,
		"currency":           "USD",
		"is_active":          true,
		"is_blocked":         false,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// GetBalance - получить баланс карты
func (h *CardHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cardID := r.URL.Query().Get("card_id")
	userID := r.URL.Query().Get("user_id")

	if cardID == "" || userID == "" {
		http.Error(w, "card_id and user_id are required", http.StatusBadRequest)
		return
	}

	h.logger.Info("fetching card balance",
		zap.String("card_id", cardID),
		zap.String("user_id", userID),
	)

	response := map[string]interface{}{
		"balance":  1500.50,
		"currency": "USD",
	}

	_ = json.NewEncoder(w).Encode(response)
}

// Deposit - пополнить баланс
func (h *CardHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		CardID      int64   `json:"card_id"`
		UserID      int64   `json:"user_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("deposit",
		zap.Int64("card_id", req.CardID),
		zap.Float64("amount", req.Amount),
	)

	response := map[string]interface{}{
		"transaction": map[string]interface{}{
			"id":               "txn-123",
			"transaction_type": "deposit",
			"amount":           req.Amount,
			"balance_after":    1500.50 + req.Amount,
			"status":           "success",
		},
		"new_balance": 1500.50 + req.Amount,
	}

	_ = json.NewEncoder(w).Encode(response)
}
