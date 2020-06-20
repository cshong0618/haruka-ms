package subscriber

import (
	"context"
	"encoding/json"
	"github.com/cshong0618/haruka/user/pkg/domain/usecase"
	"github.com/nats-io/nats.go"
)

type UserSubscriber struct {
	userService *usecase.UserService
}

func NewUserSubscriber(userService *usecase.UserService) *UserSubscriber {
	return &UserSubscriber{userService: userService}
}

type CreateUserInput struct {
	Name string `json:"name"`
}

type CreateUserOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Error *Error `json:"error"`
}

func (subscriber *UserSubscriber) CreateUser(m *nats.Msg) {
	var input CreateUserInput
	err := json.Unmarshal(m.Data, &input)
	output := CreateUserOutput{}

	if err != nil {
		output.Error = &Error{Reason: err.Error()}
		respondCreateUser(output, m)
		return
	}

	user, err := subscriber.userService.CreateUser(context.Background(), input.Name)
	if err != nil {
		output.Error = &Error{Reason: err.Error()}
		respondCreateUser(output, m)
		return
	}

	output.ID = user.ID
	output.Name = user.Name

	respondCreateUser(output, m)
}

func respondCreateUser(output CreateUserOutput, m *nats.Msg) {
	bs, _ := json.Marshal(output)
	m.Respond(bs)
}