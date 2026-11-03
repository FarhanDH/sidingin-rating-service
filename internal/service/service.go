package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/FarhanDH/sidingin-rating-service/internal/domain"
	"github.com/FarhanDH/sidingin-rating-service/internal/models"
	"github.com/FarhanDH/sidingin-rating-service/internal/repository"
)

type RatingService struct {
	repo repository.RatingRepository
}

func NewRatingService(repo repository.RatingRepository) *RatingService {
	return &RatingService{
		repo: repo,
	}
}

// TODO: Handle completed order to create new ticket rating, get active ticket by order id / ids.

func (s *RatingService) HandleCompletedOrder(ctx context.Context, payload models.OrderCompletedPayload) error {
	slog.Info("RatingService.HandleCompletedOrder",
		slog.String("order_id", payload.OrderID),
		slog.String("customer_id", payload.CustomerID),
		slog.String("technician_id", payload.TechnicianID),
		slog.Time("completed_at", payload.CompletedAt),
	)

	expiresAt := payload.CompletedAt.Add(14 * 24 * time.Hour)

	ticket := &domain.RatingTicket{
		OrderID:          payload.OrderID,
		CustomerID:       payload.CustomerID,
		TechnicianID:     payload.TechnicianID,
		ExpiresAt:        expiresAt,
		CompletedOrderAt: payload.CompletedAt,
	}

	if err := s.repo.CreateTicket(ctx, ticket); err != nil {
		slog.Error("HandleCompletedOrder: failed to create rating ticket",
			slog.String("order_id", payload.OrderID),
			slog.String("error", err.Error()),
		)
		return err
	}

	slog.Info("HandleCompletedOrder: rating ticket created successfully",
		slog.String("order_id", payload.OrderID),
		slog.Time("expires_at", expiresAt),
	)
	return nil

}

func (s *RatingService) GetBatchRatingStatus(ctx context.Context, orderIDs []string, role domain.UserRole) (map[string]string, error) {
	fmt.Printf("RatingService.GetBatchRatingStatus(), %d orderIDs, for user %s\n", len(orderIDs), string(role))
	result := make(map[string]string)

	// is orderIDs exists
	if len(orderIDs) == 0 {
		return result, nil
	}

	var targetDirection domain.DirectionRatingType
	if role == domain.CustomerRole {
		targetDirection = domain.CustomerToTechnicianDirection
	} else if role == domain.TechnicianRole {
		targetDirection = domain.TechnicianToCustomer
	} else {
		targetDirection = domain.CustomerToTechnicianDirection
	}

	// Check is the rating tickets exists based on given orderIDs
	ratingTickets, err := s.repo.GetTicketsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, err
	}

	ratedMap, err := s.repo.GetRatedOrderIDs(ctx, orderIDs, targetDirection)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	for _, ratingTicket := range ratingTickets {
		if ratedMap[ratingTicket.OrderID] {
			result[ratingTicket.OrderID] = "rated"
			continue
		}
		if now.After(ratingTicket.ExpiresAt) {
			result[ratingTicket.OrderID] = "expired"
			continue
		}

		result[ratingTicket.OrderID] = "open"
	}

	return result, nil
}
