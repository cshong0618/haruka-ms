package subscriber

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type GetUserInput struct {
	ID string `json:"id"`
}

type GetUserOutput struct {
	User
	Error *Error `json:"error"`
}

func (subscriber *UserSubscriber) GetUser(m *nats.Msg) {
	var input GetUserInput
	err := json.Unmarshal(m.Data, &input)
	output := GetUserOutput{}

	if err != nil {
		output.Error = &Error{Reason: err.Error()}
		respondGetUser(output, m)
		return
	}

	user, err := subscriber.userService.FindUserById(context.Background(), input.ID)
	if err != nil {
		output.Error = &Error{Reason: err.Error()}
		respondGetUser(output, m)
		return
	}

	output.User = toDTOUser(user)
	respondGetUser(output, m)
}

func respondGetUser(output GetUserOutput, m *nats.Msg) {
	bs, _ := json.Marshal(output)
	m.Respond(bs)
}