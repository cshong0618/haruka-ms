package main

import (
	"github.com/cshong0618/haruka/gateway/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

func main() {
	uri := os.Getenv("NATS_URL")
	port := os.Getenv("PORT")
	var err error
	var nc *nats.Conn

	for i := 0; i < 10; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		log.Printf("waiting for retry. current attempt: %d", i + 1)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		panic(err)
	}
	defer nc.Close()

	e := echo.New()

	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userHandler := user.NewUserHandler(nc)

	e.POST("/user", userHandler.CreateUser)
	e.GET("/user/:id", userHandler.FindUser)

	err = e.Start(":" + port)
	if err != nil {
		panic(err)
	}
}
