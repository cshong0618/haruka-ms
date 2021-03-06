package user

import (
	"encoding/json"
	"errors"
	"github.com/cshong0618/haruka/gateway/internal"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"time"
)

type UserHandler struct {
	nc *nats.Conn
}

func NewUserHandler(nc *nats.Conn) *UserHandler {
	return &UserHandler{nc: nc}
}

func (handler *UserHandler) CreateUser(c echo.Context) error {
	var input CreateUserInput
	if err := c.Bind(&input); err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(400, response)
		return err
	}

	bs, _ := json.Marshal(input)
	msg, err := handler.nc.Request("user.create", bs, 5*time.Second)

	if err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	var output CreateUserOutput
	err = json.Unmarshal(msg.Data, &output)
	if err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	if output.Error != nil {
		response := internal.NewErrorResponse(errors.New(output.Error.Reason))
		err = c.JSON(500, response)
		return err
	}

	// fire create auth here
	createAuthMessage := CreateAuthMessage{
		UserID:   output.ID,
		Handler:  input.Handler,
		Password: input.Password,
	}

	createAuthMessageBs, _ := json.Marshal(createAuthMessage)
	handler.nc.Publish("auth.create", createAuthMessageBs)

	user := User{
		ID:     output.ID,
		Name:   output.Name,
		Status: output.Status,
	}
	response := internal.NewResponse(user)
	err = c.JSON(200, response)
	return err
}

func (handler *UserHandler) FindUser(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		response := internal.NewResponse(errors.New("no userId provided"))
		err := c.JSON(400, response)
		return err
	}

	input := GetUserInput{ID: userID}
	bs, _ := json.Marshal(input)
	msg, err := handler.nc.Request("user.get", bs, 5*time.Second)
	if err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	var output GetUserOutput
	err = json.Unmarshal(msg.Data, &output)
	if err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	if output.Error != nil {
		response := internal.NewErrorResponse(errors.New(output.Error.Reason))
		err = c.JSON(500, response)
		return err
	}

	user := User{
		ID:     output.ID,
		Name:   output.Name,
		Status: output.Status,
	}

	response := internal.NewResponse(user)
	err = c.JSON(200, response)
	return err
}
