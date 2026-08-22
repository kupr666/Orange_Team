package domain

import (
	"time"
)

type Exercise struct {
	ID          int
	Name        string
	Description *string
	Score       int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
