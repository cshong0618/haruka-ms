package user

import "github.com/cshong0618/haruka/gateway/internal"

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Handler  string `json:"handler"`
	Password string `json:"password"`
}

type CreateAuthMessage struct {
	UserID   string `json:"userId"`
	Handler  string `json:"handler"`
	Password string `json:"password"`
}

type CreateUserOutput struct {
	User
	Error *internal.Error `json:"error"`
}

type GetUserInput struct {
	ID string `json:"id"'`
}

type GetUserOutput struct {
	User
	Error *internal.Error `json:"error"`
}
