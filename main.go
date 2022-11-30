package main

import (
	"log"
	"os"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yanamorelli/gym_consistency/database"
	"github.com/yanamorelli/gym_consistency/handlers"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	conn := os.Getenv("DBCONN")
	if conn == "" {
		log.Fatal("There isn't DBCONN variable setted.")
	}
	db, err := database.ConnectDatabase(conn)
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.Handler{DB: db}

	// TODO: Change the routes names, this isn't good
	e.POST("/", h.WentGym)
	e.GET("/getDate/:date", h.GetDate)
	e.GET("/getCurrentMonth", h.StatsOfMonth)

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.Logger.Fatal(e.Start(":8080"))
}
