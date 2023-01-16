package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yanamorelli/gym_consistency/models"
)

func ValidateGoogleJWT(tokenString string) (models.GoogleClaims, error) {
	var googleClaims models.GoogleClaims

	token, err := jwt.ParseWithClaims(tokenString, &googleClaims, func(token *jwt.Token) (interface{}, error) {
		pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, err
		}
		return key, nil
	})

	if err != nil {
		return models.GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*models.GoogleClaims)
	if !ok {
		return models.GoogleClaims{}, err
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return models.GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != "CreateClientId" {
		return models.GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return models.GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func getGooglePublicKey(keyId string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	jsonResponse := map[string]string{}
	err = json.Unmarshal(data, &jsonResponse)
	if err != nil {
		return "", nil
	}

	key, ok := jsonResponse[keyId]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}
