package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yanamorelli/gym_consistency/services"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	handler := services.SetupEnviroment()

	// TODO: Change the routes names, this isn't good
	e.POST("/", handler.WentGym)
	e.GET("/getDate/:date", handler.GetDate)
	e.GET("/getCurrentMonth", handler.StatsOfMonth)
	e.GET("/signinUser", handler.CreateUser)
	e.POST("/loginUser", handler.LoginUser)

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.Logger.Fatal(e.Start(":8080"))
}
