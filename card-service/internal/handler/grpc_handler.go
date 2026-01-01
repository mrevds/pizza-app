package handler

import (
	"card-service/internal/service"
	pbCard "card-service/pkg/card_v1/card-service_v1"
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pbCard.UnimplementedCardServiceServer
	//service слой добавить нада
	cardService service.CardService
}

func NewCardGRPCHandler(cardServ service.CardService) pbCard.CardServiceServer {
	return &grpcHandler{cardService: cardServ}
}

func (h *grpcHandler) AddCard(ctx context.Context, req *pbCard.AddCardRequest) (*pbCard.AddCardResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	isActive := true
	if req.CardData.IsActive != nil && !*req.CardData.IsActive {
		isActive = false
	}

	var dateExpire time.Time
	if req.CardData.DateExpire != nil {
		dateExpire = time.Date(
			int(req.CardData.DateExpire.Year),
			time.Month(req.CardData.DateExpire.Month),
			int(req.CardData.DateExpire.Day),
			0, 0, 0, 0, time.UTC,
		)
	}

	card, err := h.cardService.AddCard(ctx, service.AddCardData{
		Id:         uuid.NewString(),
		CardNumber: req.CardData.CardNumber,
		DateExpire: dateExpire,
		Currency:   "rub",
		Cvv:        req.Cvv,
		OwnerID:    userID,
		IsActive:   isActive,
	})

	if err != nil {
		return nil, err
	}

	return &pbCard.AddCardResponse{
		CardId:  card.ID,
		Message: "Card added successfully",
		Success: true,
	}, nil
}
