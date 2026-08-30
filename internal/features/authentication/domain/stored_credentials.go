package authentication_domain

import "github.com/google/uuid"

type StoredCredentials struct {
	UserID       uuid.UUID
	PasswordHash string
	Role         string
}
