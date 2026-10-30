package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FarhanDH/sidingin-rating-service/internal/domain"
	"github.com/FarhanDH/sidingin-rating-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository_CreateTicket(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	now := time.Now()

	// Dummy data of Rating
	ticket := &domain.RatingTicket{
		OrderID:          "order-1",
		CustomerID:       "customer-1",
		TechnicianID:     "technician-1",
		CompletedOrderAt: now,
		ExpiresAt:        now.Add(14 * 24 * time.Hour),
		ReminderCount:    0,
		LastReminderAt:   nil,
		CreatedAt:        now,
	}

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "It should successfully insert new Ticket",
			mock: func() {
				mock.ExpectExec("INSERT INTO rating_tickets").WithArgs(
					ticket.OrderID,
					ticket.CustomerID,
					ticket.TechnicianID,
					ticket.CompletedOrderAt,
					ticket.ExpiresAt,
					ticket.ReminderCount,
					ticket.LastReminderAt,
					ticket.CreatedAt,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "It should throw error when the DB is down",
			mock: func() {
				mock.ExpectExec("INSERT INTO rating_tickets").WillReturnError(errors.New("db connection failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.CreateTicket(context.Background(), ticket)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
