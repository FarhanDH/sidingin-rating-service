package models

import "time"

type OrderCompletedPayload struct {
	OrderID      string    `json:"order_id"`
	CustomerID   string    `json:"customer_id"`
	TechnicianID string    `json:"technician_id"`
	CompletedAt  time.Time `json:"completed_at"`
}
