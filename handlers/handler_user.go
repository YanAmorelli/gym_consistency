package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/yanamorelli/gym_consistency/models"
	"github.com/yanamorelli/gym_consistency/services"
)

func (h Handler) CreateUser(c echo.Context) error {
	var user models.User
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

	// TODO: Create username validation

	query := fmt.Sprintf("INSERT INTO user_info(fullname, username, passwd, email) "+
		"VALUES ('%s','%s','%s','%s') RETURNING user_id", user.FullName, user.Username, user.Password, user.Email)
	var idReturned string
	if err := h.DB.Raw(query).Scan(&idReturned).Error; err != nil {
		log.Error("error trying to create new user. Error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error": err.Error(),
		})
	}

	token, err := services.GenerateJWT(h.SecretKeyJWT, models.Claims{
		Username: user.Username,
		UserId:   idReturned,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, token)
}

func (h Handler) LoginUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		log.Error("error in json data. Error: ", err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
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

	query := fmt.Sprintf("SELECT * FROM auth_login('%s', '%s')", user.Username, user.Password)

	var userFound models.User
	if err = h.DB.Table("user_info").Raw(query).Scan(&userFound).Error; err != nil {
		log.Error("error getting user data. Error: ", err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error":  err.Error(),
			"logged": false,
		})
	}
	if userFound.Id == "" {
		return c.JSON(http.StatusNotFound, models.JsonObj{
			"error":  "username or password are incorrect",
			"logged": false,
		})
	}

	fmt.Println(userFound)

	token, err := services.GenerateJWT(h.SecretKeyJWT, models.Claims{
		Username: userFound.Username,
		UserId:   userFound.Id,
	})

	return c.JSON(http.StatusOK, models.JsonObj{
		"logged": true,
		"token":  token,
	})
}

func (h Handler) ForgetPassword(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		log.Error("error in json data. Error: ", err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	query := fmt.Sprintf("SELECT user_id,username, email FROM user_info WHERE username = '%s'", user.Username)
	err = h.DB.Raw(query).Scan(&user).Error
	if err != nil {
		log.Error("error getting user data. Error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error getting user data from database",
		})
	}

	passwordGenerated := services.GenerateRandomPassword()

	query = fmt.Sprintf("UPDATE user_info SET passwd = '%s' WHERE user_id = '%s'",
		passwordGenerated, user.Id)
	err = h.DB.Raw(query).Scan(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error updating user data from database",
		})
	}

	// TODO: Criar uma variável constante com frase e determinar idioma no futuro
	body := fmt.Sprintf("Your password was changed to %s", passwordGenerated)
	err = services.SendEmail(h.Email, h.Password, user.Email, "Forget password", body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to send email to user email",
		})
	}

	return c.JSON(http.StatusOK, models.JsonObj{
		"message": "An email was sent to " + user.Email,
	})
}

func (h Handler) ResetPassword(c echo.Context) error {
	var user models.ForgetPassword

	err := c.Bind(&user)
	if err != nil {
		log.Error("error in json data. Error: ", err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	if user.NewPassword != user.ConfirmNewPassword {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"message": "password doesn't match",
		})
	}

	query := fmt.Sprintf("SELECT * FROM auth_login('%s', '%s')", user.Username, user.OldPassword)
	var userFound models.User

	if err := h.DB.Table("user_info").Raw(query).Scan(&userFound).Error; err != nil {
		log.Error("error getting user data. Error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error": err.Error(),
		})
	}

	if userFound.Id == "" {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"message": "the actual password is different from the given password",
		})
	}

	query = fmt.Sprintf("UPDATE user_info SET passwd = '%s' WHERE username = '%s'",
		user.NewPassword, user.Username)
	if err := h.DB.Raw(query).Scan(&userFound).Error; err != nil {
		log.Error("error updating user data. Error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.JsonObj{
		"message": "user password was changed",
	})
}
