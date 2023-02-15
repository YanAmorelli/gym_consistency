package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/yanamorelli/gym_consistency/models"
	"github.com/yanamorelli/gym_consistency/services"
)

const (
	PENDENT = iota + 1
	ACCEPTED
	DECLINED
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
		(requestFriendship.RequestStatus == PENDENT || requestFriendship.RequestStatus == ACCEPTED) {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"message": fmt.Sprintf("User %s already request friendship to user %s",
				requestFriendship.UserSent, requestFriendship.UserReceived),
		})
	}

	if requestFriendship.RequestStatus == DECLINED {
		query := fmt.Sprintf(`UPDATE friend_request SET request_status=%d,dt_replied=NULL
		WHERE user_sent = '%s' AND user_received = '%s'`, PENDENT, claims.UserId, requestedUserId)
		err = h.DB.Raw(query).Scan(nil).Error
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, models.JsonObj{
				"error":   err.Error(),
				"message": "error trying to update user request",
			})
		}

		return c.JSON(http.StatusOK, nil)
	}

	query := fmt.Sprintf("INSERT INTO friend_request(user_sent, user_received, request_status) "+
		"VALUES('%s','%s',%d)", claims.UserId, requestedUserId, PENDENT)
	err = h.DB.Raw(query).Scan(&requestFriendship).Error
	if err != nil {
		log.Error(err.Error())
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
		claims.UserId, PENDENT).Find(&requestsFriendship).Error
	if err != nil {
		log.Error(err.Error())
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
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"error": err.Error(),
		})
	}

	var checkRequest models.RequestFriendship
	err = h.DB.Table("friend_request").Where("user_received = ? AND user_sent = ?", claims.UserId,
		requestFriendship.UserSent).First(&checkRequest).Error
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to get user request",
		})
	}

	if checkRequest.RequestStatus == ACCEPTED || checkRequest.RequestStatus == DECLINED {
		log.Error("Forbidden request update")
		return c.JSON(http.StatusForbidden, models.JsonObj{
			"error":   "Forbidden request update",
			"message": "Request must be pendent to be accepted ou declined",
		})
	}

	query := fmt.Sprintf(`UPDATE friend_request SET request_status = %d, dt_replied=now() WHERE
	user_sent = '%s' AND user_received = '%s' AND request_status=%d`, checkRequest.RequestStatus,
		requestFriendship.UserSent, claims.UserId, PENDENT)
	err = h.DB.Raw(query).Scan(nil).Error
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to update user request",
		})
	}

	if requestFriendship.RequestStatus == DECLINED {
		return c.JSON(http.StatusOK, nil)
	}

	insertFriendshipUserReceivedToUserSent := fmt.Sprintf(`INSERT INTO user_friendship("user", friend) 
		VALUES('%s','%s')`, claims.UserId, requestFriendship.UserSent)
	err = h.DB.Raw(insertFriendshipUserReceivedToUserSent).Scan(&requestFriendship).Error
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to insert user sent and user received in friendship table",
		})
	}

	insertFriendshipUserSentToUserReceived := fmt.Sprintf(`INSERT INTO user_friendship("user", friend)
		VALUES('%s','%s')`, requestFriendship.UserSent, claims.UserId)
	err = h.DB.Raw(insertFriendshipUserSentToUserReceived).Scan(&requestFriendship).Error
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to insert user received and user sent in friendship table",
		})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h Handler) RemoveFriend(c echo.Context) error {
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

	userTargetId := c.Param("userId")
	query := fmt.Sprintf(`DELETE FROM user_friendship WHERE "user"='%s' AND friend='%s' OR 
	"user"='%s' AND friend='%s'`, claims.UserId, userTargetId, userTargetId, claims.UserId)

	err = h.DB.Raw(query).Scan(nil).Error
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"message": "error trying to delete user",
			"error":   err.Error(),
		})
	}

	err = h.DB.Table("friend_request").Where("user_sent = ? AND user_received = ? OR user_sent = ? AND user_received = ?",
		claims.UserId, userTargetId, userTargetId, claims.UserId).Update("request_status", DECLINED).Error
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error":   err.Error(),
			"message": "error trying to update user request",
		})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h Handler) GetUserFriends(c echo.Context) error {

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

	var friends []string

	err = h.DB.Table("user_friendship").Select("friend").Where(`"user"=?`, claims.UserId).Find(&friends).Error
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"message": "error trying to get user data",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.JsonObj{
		"friends": friends,
	})
}
