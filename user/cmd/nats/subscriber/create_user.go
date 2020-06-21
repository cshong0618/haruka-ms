package subscriber

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type CreateUserInput struct {
	Name string `json:"name"`
}

type CreateUserOutput struct {
	User
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

	output.User = toDTOUser(user)
	respondCreateUser(output, m)
}

func respondCreateUser(output CreateUserOutput, m *nats.Msg) {
	bs, _ := json.Marshal(output)
	m.Respond(bs)
}