package jwt_services

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	"github.com/yanamorelli/gym_consistency/models"
)

func GenerateJWT(secretKey string, claims models.Claims) (string, error) {
	// signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secretKey))
	// if err != nil {
	// 	return "", err
	// }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println(err.Error())
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
		log.Error(err)
		return claims, err
	}
	if !token.Valid {
		log.Error("error in validation")
		return claims, errors.New("invalid token")
	}

	return claims, nil
}
