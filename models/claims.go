package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	jwt.RegisteredClaims
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}
