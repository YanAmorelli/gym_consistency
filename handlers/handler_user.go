package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yanamorelli/gym_consistency/models"
	"github.com/yanamorelli/gym_consistency/models/model_user"
	"github.com/yanamorelli/gym_consistency/services/jwt_services"
)

func (h Handler) CreateUser(c echo.Context) error {
	var user model_user.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	token, err := jwt_services.GenerateJWT(h.SecretKeyJWT, models.Claims{})

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, token)
}

func (h Handler) LoginUser(c echo.Context) error {
	token := c.Request().Header.Get("Token")
	fmt.Println(token)

	claims, err := jwt_services.VerifyJWT(token, h.SecretKeyJWT)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error": err.Error(),
		})
	}

	fmt.Println(claims)

	return c.JSON(http.StatusOK, nil)
}
