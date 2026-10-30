package domain

import "time"

type UserRole string

const (
	CustomerRole   UserRole = "customer"
	TechnicianRole UserRole = "technician"
)

type DirectionRatingType string

const (
	CustomerToTechnicianDirection DirectionRatingType = "customer_to_technician"
	CustomerToServiceDirection    DirectionRatingType = "customer_to_service"
	CustomerToAppDirection        DirectionRatingType = "customer_to_app"
	TechnicianToCustomer          DirectionRatingType = "technician_to_customer"
	TechnicianToApp               DirectionRatingType = "technician_to_app"
)

type RatingTicket struct {
	OrderID          string
	CustomerID       string
	TechnicianID     string
	CompletedOrderAt time.Time
	ExpiresAt        time.Time
	ReminderCount    int
	// It can be nil if it has never been reminded
	LastReminderAt *time.Time
	CreatedAt      time.Time
}

type RatingTag struct {
	Code      string
	Direction DirectionRatingType
	Label     string
	MinStars  int
	MaxStars  int
	SortOrder int
	IsActive  bool
}

type Rating struct {
	ID               string
	OrderID          string
	Direction        DirectionRatingType
	Stars            int
	Comment          *string
	IsCommentHidden  bool
	IsRequiresReview bool
	// AppPlatform is nil when the direction rating is not to App
	AppPlatform *string
	// AppVersion is nil when the direction rating is not to App
	AppVersion *string
	CreatedAt  time.Time
}

type RatingSelection struct {
	RatingID      string
	Direction     DirectionRatingType
	RatingTagCode string
}
