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
	if SECRET_KEY_JWT == "" {
		log.Fatal("There isn't SECRET_KEY_JWT variable setted")
	}

	EMAIL_GYM := os.Getenv("EMAIL_GYM")
	if EMAIL_GYM == "" {
		log.Fatal("There isn't EMAIL_GYM variable setted")
	}

	PASSWORD_EMAIL_GYM := os.Getenv("PASSWORD_EMAIL_GYM")
	if PASSWORD_EMAIL_GYM == "" {
		log.Fatal("There isn't PASSWORD_EMAIL_GYM variable setted")
	}

	db, err := database.ConnectDatabase(conn)
	if err != nil {
		log.Fatal(err)
	}

	return handlers.Handler{
		DB:           db,
		SecretKeyJWT: SECRET_KEY_JWT,
		Email:        EMAIL_GYM,
		Password:     PASSWORD_EMAIL_GYM,
	}
}
