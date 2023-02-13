package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/yanamorelli/gym_consistency/models"
	"github.com/yanamorelli/gym_consistency/services"
)

func (h Handler) RequestFriendship(c echo.Context) error {
	token := c.Request().Header.Get("token")
	if token == "" {
		message := "token not provided"
		log.Error(message)
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": message,
		})
	}
	claims, err := services.VerifyJWT(token, h.SecretKeyJWT)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error":   err.Error(),
			"message": "error in token verification",
		})
	}

	requestedUserId := c.Param("requestedUserId")
	var requestFriendship models.RequestFriendship

	err = h.DB.Table("friend_request").Select("request_id, user_sent, user_received, dt_sented, dt_replied, request_status").
		Where("user_sent=? AND user_received=?", claims.UserId, requestedUserId).
		First(&requestFriendship).Error

	if err != nil && err.Error() != "record not found" {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error seeking for friendship request",
		})
	}

	if requestFriendship.UserSent == claims.UserId &&
		requestFriendship.UserReceived == requestedUserId &&
		(requestFriendship.RequestStatus == 1 || requestFriendship.RequestStatus == 2) {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"message": fmt.Sprintf("User %s already request friendship to user %s",
				requestFriendship.UserSent, requestFriendship.UserReceived),
		})
	}

	if requestFriendship.RequestStatus == 3 {
		err = h.DB.Table("friend_request").Where("user_sent = ? AND user_received = ?", claims.UserId, requestedUserId).
			Update("request_status", 1).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JsonObj{
				"error":   err.Error(),
				"message": "error trying to update user request",
			})
		}

		return c.JSON(http.StatusOK, nil)
	}

	query := fmt.Sprintf("INSERT INTO friend_request(user_sent, user_received, request_status) "+
		"VALUES('%s','%s',%d)", claims.UserId, requestedUserId, 1)
	err = h.DB.Raw(query).Scan(&requestFriendship).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to create friendship user request",
		})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h Handler) GetFriendshipRequest(c echo.Context) error {
	token := c.Request().Header.Get("token")
	if token == "" {
		message := "token not provided"
		log.Error(message)
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": message,
		})
	}
	claims, err := services.VerifyJWT(token, h.SecretKeyJWT)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error":   err.Error(),
			"message": "error in token verification",
		})
	}

	var requestsFriendship []models.RequestFriendship

	err = h.DB.Table("friend_request").Where("user_received = ? AND request_status = ?",
		claims.UserId, 1).Find(&requestsFriendship).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to get friendship requests",
		})
	}

	return c.JSON(http.StatusOK, requestsFriendship)
}

func (h Handler) UpdateFriendshipRequest(c echo.Context) error {
	token := c.Request().Header.Get("token")
	if token == "" {
		message := "token not provided"
		log.Error(message)
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": message,
		})
	}
	claims, err := services.VerifyJWT(token, h.SecretKeyJWT)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error":   err.Error(),
			"message": "error in token verification",
		})
	}

	var requestFriendship models.RequestFriendship
	err = c.Bind(&requestFriendship)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	var checkRequest models.RequestFriendship
	err = h.DB.Table("friend_request").Where("user_received = ? AND user_sent = ?", claims.UserId,
		requestFriendship.UserSent).First(&checkRequest).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to get user request",
		})
	}

	if checkRequest.RequestStatus == 2 || checkRequest.RequestStatus == 3 {
		return c.JSON(http.StatusForbidden, models.JsonObj{
			"error":   "Forbidden request update",
			"message": "Request must be pendent to be accepted ou declined",
		})
	}

	err = h.DB.Table("friend_request").Where("user_sent = ? AND user_received = ? AND request_status=1",
		requestFriendship.UserSent, claims.UserId).Update("request_status", requestFriendship.RequestStatus).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to update user request",
		})
	}

	if requestFriendship.RequestStatus == 3 {
		return c.JSON(http.StatusOK, nil)
	}

	insertFriendshipUserReceivedToUserSent := fmt.Sprintf(`INSERT INTO user_friendship("user", friend) 
		VALUES('%s','%s')`, claims.UserId, requestFriendship.UserSent)
	err = h.DB.Raw(insertFriendshipUserReceivedToUserSent).Scan(&requestFriendship).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to insert user sent and user received in friendship table",
		})
	}

	insertFriendshipUserSentToUserReceived := fmt.Sprintf(`INSERT INTO user_friendship("user", friend)
		VALUES('%s','%s')`, requestFriendship.UserSent, claims.UserId)
	err = h.DB.Raw(insertFriendshipUserSentToUserReceived).Scan(&requestFriendship).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to insert user received and user sent in friendship table",
		})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h Handler) RemoveFriend(c echo.Context) error {
	// Validar token do usuário
	// Remove usuário amizade dos usuários
	return nil
}

func (h Handler) GetUserFriends(c echo.Context) error {
	// Validar token do usuário
	// Pegar id do usuário no token
	// Buscar amigos do usuário
	// Retornar lista de amigos
	return nil
}
