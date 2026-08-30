package core_auth_jwt

import "github.com/golang-jwt/jwt/v5"

type claims struct {
	Role string `json:"role,omitempty"`
	jwt.RegisteredClaims
}
