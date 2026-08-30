package core_auth

import "github.com/google/uuid"

type Principal struct {
	UserID uuid.UUID
	Role   string
}
