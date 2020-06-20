package main

import (
	"github.com/cshong0618/haruka/user/cmd/rest/controller"
	"github.com/cshong0618/haruka/user/cmd/wire"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	// Manual DI
	userService := wire.InitUserService()

	// Server setup
	e := echo.New()
	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userController := controller.NewUserController(userService)
	e.POST("/user", userController.CreateUser)

	err := e.Start(":" + port)
	if err != nil {
		panic(err)
	}
}
