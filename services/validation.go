package services

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/yanamorelli/gym_consistency/models"
)

func ValidateUserData(user *models.User) error {
	if user.Username == "" {
		return errors.New("username not provided")
	}

	if user.Email == "" {
		return errors.New("email not provided")
	}

	if user.FullName == "" {
		return errors.New("fullname not provided")
	}

	if user.Password == "" {
		return errors.New("password not provided")
	}

	if user.ConfirmPassword == "" {
		return errors.New("confirmation password not provided")
	}

	if user.ConfirmPassword != user.Password {
		return errors.New("passwords doesn't match")
	}

	return nil
}

func ValidateUserToken(c echo.Context, secretKey string) (models.Claims, error) {
	token, err := c.Request().Cookie("token")
	if err != nil {
		log.Error(err.Error())
		return models.Claims{}, err
	}
	if token.Value == "" {
		message := "token not provided"
		log.Error(message)
		return models.Claims{}, errors.New(message)
	}
	claims, err := VerifyJWT(token.Value, secretKey)
	if err != nil {
		log.Error(err.Error())
		return models.Claims{}, err
	}
	return *claims, nil
}
