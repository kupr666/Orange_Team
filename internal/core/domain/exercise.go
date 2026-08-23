package domain

import (
	"time"

	"github.com/google/uuid"
)

type Exercise struct {
	ID          uuid.UUID
	Version     int
	Name        string
	Description string
	Difficulty  int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Type        string
}

func NewExercise(
	id uuid.UUID,
	version int,
	name string,
	description string,
	difficulty int,
	createdAt time.Time,
	updatedAt *time.Time,
	exerciseType string,
) Exercise {
	return Exercise{
		ID:          id,
		Version:     version,
		Name:        name,
		Description: description,
		Difficulty:  difficulty,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Type:        exerciseType,
	}
}
