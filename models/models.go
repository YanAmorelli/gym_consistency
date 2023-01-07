package models

import "github.com/golang-jwt/jwt/v4"

type Ok struct {
	Date string `json:"date" gorm:"column:date_gym"`
	Ok   bool   `json:"ok" gorm:"column:ok"`
}

type Stats struct {
	PresentDays int64 `json:"presentDays"`
	MissedDays  int64 `json:"missedDays"`
}

type JsonObj map[string]interface{}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
