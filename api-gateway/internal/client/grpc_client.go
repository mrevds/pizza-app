package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	UserServiceConn *grpc.ClientConn
	CardServiceConn *grpc.ClientConn
}

func NewGRPCClients(userServiceAddr, cardServiceAddr string) (*GRPCClients, error) {
	// User Service
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	// Card Service
	cardConn, err := grpc.Dial(cardServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = userConn.Close()
		return nil, fmt.Errorf("failed to connect to card service: %w", err)
	}

	return &GRPCClients{
		UserServiceConn: userConn,
		CardServiceConn: cardConn,
	}, nil
}

func (c *GRPCClients) Close() error {
	if err := c.UserServiceConn.Close(); err != nil {
		return fmt.Errorf("failed to close user service connection: %w", err)
	}
	if err := c.CardServiceConn.Close(); err != nil {
		return fmt.Errorf("failed to close card service connection: %w", err)
	}
	return nil
}
