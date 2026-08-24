package authentication_service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

type authenticationRepositoryStub struct {
	createdEmail        string
	createdPasswordHash string
	createdFullName     string
	err                 error
}

func (r *authenticationRepositoryStub) CreateUser(
	_ context.Context,
	email string,
	passwordHash string,
	fullName string,
) (domain.User, error) {
	r.createdEmail = email
	r.createdPasswordHash = passwordHash
	r.createdFullName = fullName
	if r.err != nil {
		return domain.User{}, r.err
	}

	return domain.User{
		ID:       uuid.New(),
		Email:    email,
		FullName: fullName,
	}, nil
}

func TestRegisterUserNormalizesEmailAndHashesPassword(t *testing.T) {
	repository := &authenticationRepositoryStub{}
	service := NewAuthenticationService(repository)
	plainPassword := "secret123"

	createdUser, err := service.RegisterUser(
		context.Background(),
		"  IVAN@example.COM  ",
		plainPassword,
		"Ivan Ivanov",
	)
	if err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}

	if repository.createdEmail != "ivan@example.com" {
		t.Fatalf("repository email = %q, want normalized email", repository.createdEmail)
	}
	if repository.createdPasswordHash == plainPassword {
		t.Fatal("repository received plaintext password")
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(repository.createdPasswordHash),
		[]byte(plainPassword),
	); err != nil {
		t.Fatalf("repository received invalid password hash: %v", err)
	}
	if createdUser.Email != "ivan@example.com" {
		t.Fatalf("created user email = %q", createdUser.Email)
	}
}

func TestRegisterUserRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		fullName string
	}{
		{name: "invalid email", email: "invalid", password: "secret123", fullName: "Ivan Ivanov"},
		{name: "short password", email: "ivan@example.com", password: "short", fullName: "Ivan Ivanov"},
		{name: "password over bcrypt limit", email: "ivan@example.com", password: strings.Repeat("a", 73), fullName: "Ivan Ivanov"},
		{name: "invalid full name", email: "ivan@example.com", password: "secret123", fullName: " Ivan Ivanov"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &authenticationRepositoryStub{}
			service := NewAuthenticationService(repository)

			_, err := service.RegisterUser(
				context.Background(),
				test.email,
				test.password,
				test.fullName,
			)
			if !errors.Is(err, core_errors.ErrInvalidArgument) {
				t.Fatalf("RegisterUser() error = %v, want ErrInvalidArgument", err)
			}
			if repository.createdEmail != "" {
				t.Fatal("repository was called for invalid input")
			}
		})
	}
}

func TestRegisterUserPreservesRepositoryConflict(t *testing.T) {
	repository := &authenticationRepositoryStub{err: core_errors.ErrConflict}
	service := NewAuthenticationService(repository)

	_, err := service.RegisterUser(
		context.Background(),
		"ivan@example.com",
		"secret123",
		"Ivan Ivanov",
	)
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("RegisterUser() error = %v, want ErrConflict", err)
	}
}
