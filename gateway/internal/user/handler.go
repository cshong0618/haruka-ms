package user

import (
	"encoding/json"
	"errors"
	"github.com/cshong0618/haruka/gateway/internal"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"time"
)

type UserHandler struct{
	nc *nats.Conn
}

func NewUserHandler(nc *nats.Conn) *UserHandler {
	return &UserHandler{nc: nc}
}

func (handler *UserHandler) CreateUser(c echo.Context) error {
	var input CreateUserInput
	if err := c.Bind(&input); err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	bs, _ := json.Marshal(input)
	msg, err := handler.nc.Request("user.create", bs, 5 * time.Second)

	if err != nil {
		response := internal.NewErrorResponse(err)
		err = c.JSON(500, response)
		return err
	}

	var output CreateUserOutput
	json.Unmarshal(msg.Data, &output)

	if output.Error != nil {
		response := internal.NewErrorResponse(errors.New(output.Error.Reason))
		err = c.JSON(500, response)
		return err
	}

	user := User{
		ID:   output.ID,
		Name: output.Name,
	}
	response := internal.NewResponse(user)
	err = c.JSON(200, response)
	return err
}
