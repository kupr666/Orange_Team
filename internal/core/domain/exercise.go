package domain

import (
	"time"
)

type Exercise struct {
	ID          int
	Version     int64
	Name        string
	Description string
	Difficulty  int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Type        string
}

func NewExercise(
	id int,
	version int64,
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
