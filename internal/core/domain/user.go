package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

const (
	MinEmailLength = 5
	MaxEmailLength = 30

	MinFullNameLength = 2
	MaxFullNameLength = 50

	MinUserWorkoutScore = 0

	SexMale   = "male"
	SexFemale = "female"

	MinWeightGrams = 20000
	MaxWeightGrams = 300000

	MinHeightCM = 100
	MaxHeightCM = 250

	MinBirthDateYear = 1900
)

var EmailPattern = regexp.MustCompile(`(?i)^[a-z0-9]([a-z0-9]|[.](?![.]))*[a-z0-9]@[a-z0-9.-]+\.[a-z]{2,}$`)

var AllowedSexes = map[string]bool{
	SexMale:   true,
	SexFemale: true,
}

type User struct {
	ID               uuid.UUID
	Version          int
	Email            string
	FullName         string
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	UserWorkoutScore int
	Sex              *string
	WeightGrams      *int
	BirthDate        *time.Time
	HeightCM         *int
}

func NewUser(
	id uuid.UUID,
	version int,
	email string,
	fullName string,
	createdAt time.Time,
	updatedAt *time.Time,
	userWorkoutScore int,
	sex *string,
	weightGrams *int,
	birthDate *time.Time,
	heightCM *int,
) User {
	return User{
		ID:               id,
		Version:          version,
		Email:            email,
		FullName:         fullName,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		UserWorkoutScore: userWorkoutScore,
		Sex:              sex,
		WeightGrams:      weightGrams,
		BirthDate:        birthDate,
		HeightCM:         heightCM,
	}
}

func CreateUser(
	email string,
	fullName string,
) User {
	var (
		id                          = uuid.New()
		version                     = 1
		createdAt                   = time.Now()
		updatedAt        *time.Time = nil
		userWorkoutScore            = 0
		sex              *string    = nil
		weightGrams      *int       = nil
		birthDate        *time.Time = nil
		heightCM         *int       = nil
	)

	return NewUser(
		id,
		version,
		email,
		fullName,
		createdAt,
		updatedAt,
		userWorkoutScore,
		sex,
		weightGrams,
		birthDate,
		heightCM,
	)
}

func (u *User) Validate() error {
	if len(u.Email) < 5 || len(u.Email) > 30 {
		return fmt.Errorf(
			"email length must be between %d and %d: %w",
			MinEmailLength, MaxEmailLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if !EmailPattern.MatchString(u.Email) {
		return fmt.Errorf(
			"invalid email format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	trimmedFullName := strings.TrimSpace(u.FullName)
	if u.FullName != trimmedFullName {
		return fmt.Errorf(
			"full_name cannot have leading/trailing spaces: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	fullNameLen := len([]rune(trimmedFullName))
	if fullNameLen < MinFullNameLength || fullNameLen > MaxFullNameLength {
		return fmt.Errorf(
			"full_name length must be between %d and %d: %w",
			MinFullNameLength, MaxFullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.UserWorkoutScore < MinUserWorkoutScore {
		return fmt.Errorf(
			"user_workout_score must be >= %d: %w",
			MinUserWorkoutScore,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.Sex != nil {
		if !AllowedSexes[*u.Sex] {
			return fmt.Errorf(
				"invalid sex: %s (allowed: %s, %s): %w",
				*u.Sex, SexMale, SexFemale,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if u.WeightGrams != nil {
		if *u.WeightGrams < MinWeightGrams || *u.WeightGrams > MaxWeightGrams {
			return fmt.Errorf(
				"weight_grams must be between %d and %d: %w",
				MinWeightGrams, MaxWeightGrams,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if u.BirthDate != nil {
		minDate := time.Date(MinBirthDateYear, 1, 1, 0, 0, 0, 0, time.UTC)
		if u.BirthDate.Before(minDate) || u.BirthDate.After(time.Now()) {
			return fmt.Errorf(
				"birth_date must be between %d-01-01 and today: %w",
				MinBirthDateYear,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if u.HeightCM != nil {
		if *u.HeightCM < MinHeightCM || *u.HeightCM > MaxHeightCM {
			return fmt.Errorf(
				"height_cm must be between %d and %d: %w",
				MinHeightCM, MaxHeightCM,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

// Возможно стоит удалить
// func (u User) ProfileCompleted() bool {
// 	return u.Sex != nil &&
// 		u.WeightGrams != nil &&
// 		u.BirthDate != nil &&
// 		u.HeightCM != nil
// }

type UserPatch struct {
	Sex         Nullable[string]
	WeightGrams Nullable[int]
	BirthDate   Nullable[time.Time]
	HeightCM    Nullable[int]
}

func NewUserPatch(
	sex Nullable[string],
	weightGrams Nullable[int],
	birthDate Nullable[time.Time],
	heightCM Nullable[int],
) UserPatch {
	return UserPatch{
		Sex:         sex,
		WeightGrams: weightGrams,
		BirthDate:   birthDate,
		HeightCM:    heightCM,
	}
}

func (p *UserPatch) Validate() error {
	// Логика для валидации патча. Пока ничего нет. Может и не будет
	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf(
			"validate user patch: %w",
			err,
		)
	}

	tmp := *u

	if patch.Sex.Set {
		tmp.Sex = patch.Sex.Value
	}
	if patch.WeightGrams.Set {
		tmp.WeightGrams = patch.WeightGrams.Value
	}
	if patch.BirthDate.Set {
		tmp.BirthDate = patch.BirthDate.Value
	}
	if patch.HeightCM.Set {
		tmp.HeightCM = patch.HeightCM.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf(
			"validate patched user: %w",
			err,
		)
	}

	*u = tmp
	return nil
}
