package user

import "github.com/cshong0618/haruka/gateway/internal"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUserInput struct {
	Name string `json:"name"`
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