package subscriber

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type ActivateUserInput struct {
	ID string `json:"id"`
}

type ActivateUserOutput struct {
	User
	Error *Error `json:"error"`
}

func (subscriber *UserSubscriber) ActivateUser(m *nats.Msg) {
	var input ActivateUserInput
	err := json.Unmarshal(m.Data, &input)
	output := ActivateUserOutput{}
	if err != nil {
		output.Error = &Error{Reason: err.Error()}
		respondActivateUser(output, m)
		return
	}

	user, err := subscriber.userService.ActivateUser(context.Background(), input.ID)
	if err != nil {
		output.Error = &Error{Reason: err.Error()}
		respondActivateUser(output, m)
		return
	}

	output.User = toDTOUser(user)
	respondActivateUser(output, m)
}

func respondActivateUser(output ActivateUserOutput, m *nats.Msg) {
	bs, _ := json.Marshal(output)
	m.Respond(bs)
}