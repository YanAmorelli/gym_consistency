package setup

import (
	"log"
	"os"

	"github.com/yanamorelli/gym_consistency/database"
	"github.com/yanamorelli/gym_consistency/handlers"
)

func SetupEnviroment() handlers.Handler {
	conn := os.Getenv("DBCONN")
	if conn == "" {
		log.Fatal("There isn't DBCONN variable setted.")
	}

	SECRET_KEY_JWT := os.Getenv("SECRET_KEY_JWT")
	if conn == "" {
		log.Fatal("There isn't SECRET_KEY_JWT variable setted")
	}

	db, err := database.ConnectDatabase(conn)
	if err != nil {
		log.Fatal(err)
	}

	return handlers.Handler{
		DB:           db,
		SecretKeyJWT: SECRET_KEY_JWT,
	}
}
