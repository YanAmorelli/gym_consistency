package handlers

import (
	"gorm.io/gorm"
)

type Handler struct {
	DB           *gorm.DB
	SecretKeyJWT string
	Email        string
	Password     string
}
