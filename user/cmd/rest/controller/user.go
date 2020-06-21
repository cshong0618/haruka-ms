package controller

import (
	"github.com/cshong0618/haruka/user/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService *usecase.UserService
}

func NewUserController(userService *usecase.UserService) *UserController {
	return &UserController{userService: userService}
}


type CreateUserInput struct {
	Name string `json:"name"`
}

type CreateUserOutput struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func (controller *UserController) CreateUser(c echo.Context) error {
	var input CreateUserInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	user, err := controller.userService.CreateUser(c.Request().Context(), input.Name)

	if err != nil {
		response := NewErrorResponse(err)
		err = c.JSON(200, response)
		return err
	}

	output := CreateUserOutput{
		ID:   user.ID,
		Name: user.Name,
	}

	response := NewResponse(output)
	err = c.JSON(200, response)
	return err
}
