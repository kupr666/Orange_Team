package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	Version          int64
	Email            string
	FullName         string
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	UserWorkoutScore int
	Sex              string
	WeightGrams      *int
	BirthDate        *time.Time
	HeightCM         *int
}

func (u User) ProfileCompleted() bool {
	return u.WeightGrams != nil &&
		u.BirthDate != nil &&
		u.HeightCM != nil
}
