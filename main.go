package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yanamorelli/gym_consistency/setup"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	handler := setup.SetupEnviroment()

	e.POST("/bookAttendance", handler.WentGym)
	e.GET("/getDate/:date", handler.GetDate)
	e.GET("/getCurrentMonth", handler.StatsOfMonth)
	e.POST("/signUpUser", handler.CreateUser)
	e.POST("/loginUser", handler.LoginUser)
	e.POST("/forgetPassword", handler.ForgetPassword)
	e.POST("/resetPassword", handler.ResetPassword)
	e.POST("/requestFriendship/:requestedUserId", handler.RequestFriendship)
	e.GET("/getUserFriendshipRequests", handler.GetFriendshipRequest)
	e.POST("/updateFriendshipRequest", handler.UpdateFriendshipRequest)

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.Logger.Fatal(e.Start(":8080"))
}
