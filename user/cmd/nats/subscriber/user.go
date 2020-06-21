package subscriber

import (
	userDomain "github.com/cshong0618/haruka/user/pkg/domain/user"
	"github.com/cshong0618/haruka/user/pkg/usecase"
)

type UserSubscriber struct {
	userService *usecase.UserService
}

func NewUserSubscriber(userService *usecase.UserService) *UserSubscriber {
	return &UserSubscriber{userService: userService}
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func toDTOUser(user userDomain.User) User {
	return User{
		ID:     user.ID,
		Name:   user.Name,
		Status: string(user.UserStatus),
	}
}
