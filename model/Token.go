package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// Token model
type Token struct {
	UserID string
	jwt.StandardClaims
}
