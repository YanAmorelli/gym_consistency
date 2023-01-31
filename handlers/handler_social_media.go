package handlers

import "github.com/labstack/echo/v4"

func (h Handler) RequestFriendship(c echo.Context) error {
	// Validar token usuário
	// Receber  id do usuário solicitante pelo token e id do usuário solicitado por json ou parametro
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
