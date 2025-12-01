package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64      `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"` // null, если нет
}
