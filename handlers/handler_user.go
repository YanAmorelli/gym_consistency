package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/yanamorelli/gym_consistency/models"
	"github.com/yanamorelli/gym_consistency/models/model_user"
	"github.com/yanamorelli/gym_consistency/services"
)

func (h Handler) CreateUser(c echo.Context) error {
	var user model_user.User
	err := c.Bind(&user)
	if err != nil {
		log.Error("error in json data. Error: ", err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	validData := services.ValidateUserData(&user)
	if validData != nil {
		log.Error("error in user validation. Error: ", validData.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": validData.Error(),
		})
	}

	if err := h.DB.Table("user_info").Select("user_id", "fullname", "username", "passwd", "email").
		Create(user).Error; err != nil {
		log.Error("error trying to create new user. Error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error": err.Error(),
		})
	}

	token, err := services.GenerateJWT(h.SecretKeyJWT, models.Claims{
		Username: user.Username,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, token)
}

func (h Handler) LoginUser(c echo.Context) error {
	// TODO: Increase login send data from token
	var user model_user.User
	err := c.Bind(&user)
	if err != nil {
		log.Error("error in json data. Error: ", err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error":  err.Error(),
			"logged": false,
		})
	}
	token := c.Request().Header.Get("Token")

	_, err = services.VerifyJWT(token, h.SecretKeyJWT)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":  err.Error(),
			"logged": false,
		})
	}

	if user.Password == "" || user.Username == "" {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error":  "username and password must be provide",
			"logged": false,
		})
	}

	query := fmt.Sprintf("SELECT username FROM login_username('%s', '%s')", user.Username, user.Password)
	var userChecked model_user.User
	if err := h.DB.Table("user_info").Raw(query).Scan(&userChecked).Error; err != nil {
		log.Error("error getting user data. Error: ", err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error":  err.Error(),
			"logged": false,
		})
	}

	if userChecked.Username == "" {
		return c.JSON(http.StatusNotFound, models.JsonObj{
			"error":  "username or password are incorrect",
			"logged": false,
		})
	}

	return c.JSON(http.StatusOK, models.JsonObj{
		"logged": true,
	})
}
