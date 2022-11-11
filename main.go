package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/yanamorelli/gym_consistency/database"
	"github.com/yanamorelli/gym_consistency/handlers"
)

func main() {
	e := echo.New()

	conn := os.Getenv("DBCONN")
	if conn == "" {
		log.Fatal("There isn't DBCONN variable setted.")
	}
	db, err := database.ConnectDatabase(conn)
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.Handler{DB: db}

	e.POST("/", h.WentGym)
	e.GET("/getDate/:date", h.GetDate)
	e.GET("/getCurrentMonth", h.StatsOfMonth)

	e.Logger.Fatal(e.Start(":8080"))
}
