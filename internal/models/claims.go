package models

import (
	"github.com/golang-jwt/jwt"
)

// Claims представляет пользовательские утверждения JWT.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
