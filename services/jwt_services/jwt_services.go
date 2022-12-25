package jwt_services

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	"github.com/yanamorelli/gym_consistency/models"
)

func GenerateJWT(secretKey string, claims models.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Error("error in token generation. Error: ", err.Error())
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string, secretKey string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Error("error trying to parse token. Error: ", err.Error())
		return claims, err
	}
	if !token.Valid {
		log.Error("error checking if token is valid")
		return claims, errors.New("invalid token")
	}

	return claims, nil
}
