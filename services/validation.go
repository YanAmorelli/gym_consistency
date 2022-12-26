package services

import (
	"errors"

	"github.com/yanamorelli/gym_consistency/models/model_user"
)

func ValidateUserData(user *model_user.User) error {
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
