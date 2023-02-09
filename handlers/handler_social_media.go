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
		requestFriendship.UserReceived == requestedUserId {
		return c.JSON(http.StatusBadRequest, models.JsonObj{
			"message": fmt.Sprintf("User %s already request friendship to user %s",
				requestFriendship.UserSent, requestFriendship.UserReceived),
		})
	}

	fmt.Println("request friendship", requestFriendship)

	// Verificar se registro existe ou se status é igual de aceito, caso sim, retorna um erro de que requisiçaõ já foi feita
	// Salva requisição no banco de dados
	// Retorna validado
	return nil
}

func (h Handler) GetFriendshipRequest(c echo.Context) error {
	// Validar token do usuário
	// Receber id do usuário
	// Buscar requisições pendentes
	// Retornas requisições do usuário
	return nil
}

func (h Handler) UpdateFriendshipRequest(c echo.Context) error {
	// Validar token  do usuário
	// Atualiza requisição de amizade de acordo com json recebido
	return nil
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
