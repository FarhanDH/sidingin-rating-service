package repository

import (
	"context"
	"fmt"

	"github.com/FarhanDH/sidingin-rating-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) RatingRepository {
	return &postgresRepo{
		db: db,
	}
}

func (r *postgresRepo) CreateTicket(ctx context.Context, ticket *domain.RatingTicket) error {
	fmt.Printf("inserting ticket rating into db: %v", *ticket)
	query := `
		INSERT INTO rating_tickets (
			order_id,
			customer_id,
			technician_id,
			completed_order_at,
			expires_at,
			reminder_count,
			last_reminder_at,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		ticket.OrderID,
		ticket.CustomerID,
		ticket.TechnicianID,
		ticket.CompletedOrderAt,
		ticket.ExpiresAt,
		ticket.ReminderCount,
		ticket.LastReminderAt,
		ticket.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert rating ticket: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetTicketsByOrderIDs(ctx context.Context, orderIDs []string) ([]domain.RatingTicket, error) {
	fmt.Printf("GetTicketsByOrderIDs(), for %d orders", len(orderIDs))
	query := `
		SELECT order_id, expires_at
		FROM rating_tickets
		WHERE order_id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, pq.Array(orderIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to query rating tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]domain.RatingTicket, 0)
	for rows.Next() {
		var ticket domain.RatingTicket
		if err := rows.Scan(&ticket.OrderID, &ticket.ExpiresAt); err != nil {
			return nil, fmt.Errorf("failed to scan rating ticket row: %w", err)
		}
		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during iterate rating tickets: %w", err)
	}

	return tickets, nil
}

func (r *postgresRepo) GetRatedOrderIDs(ctx context.Context, orderIDs []string, direction domain.DirectionRatingType) (map[string]bool, error) {
	fmt.Printf("GetRatedOrderIDs(), for %d orders with direction: %v", len(orderIDs), direction)

	query := `
		SELECT order_id
		FROM ratings
		WHERE order_id = ANY($1) AND direction = $2
	`

	rows, err := r.db.Query(ctx, query, pq.Array(orderIDs), string(direction))
	if err != nil {
		return nil, fmt.Errorf("failed to query rated order ids: %w", err)
	}
	defer rows.Close()

	ratedMap := make(map[string]bool)
	for rows.Next() {
		var rating domain.Rating
		if err := rows.Scan(&rating.OrderID); err != nil {
			return nil, fmt.Errorf("failed  to scan rating row: %w", err)
		}
		ratedMap[rating.OrderID] = true
	}
	return ratedMap, nil
}
