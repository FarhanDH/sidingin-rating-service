package repository

import (
	"context"

	"github.com/FarhanDH/sidingin-rating-service/internal/domain"
)

type RatingRepository interface {
	CreateTicket(ctx context.Context, ticket *domain.RatingTicket) error
	GetTicketsByOrderIDs(ctx context.Context, orderIDs []string) ([]domain.RatingTicket, error)
	GetRatedOrderIDs(ctx context.Context, orderIDs []string, direction domain.DirectionRatingType) (map[string]bool, error)
}
