package subscriber

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type DeactivateUserInput struct {
	ID string `json:"id"`
}

type DeactivateUserOutput struct {
	OK    bool   `json:"ok"`
	Error *Error `json:"error"`
}

func (subscriber *UserSubscriber) DeactivateUser(m *nats.Msg) {
	var input DeactivateUserInput
	err := json.Unmarshal(m.Data, &input)
	output := DeactivateUserOutput{}
	if err != nil {
		output.OK = false
		output.Error = &Error{Reason: err.Error()}
		respondDeactivateUser(output, m)
		return
	}

	_, err = subscriber.userService.DeactivateUser(context.Background(), input.ID)
	if err != nil {
		output.OK = false
		output.Error = &Error{Reason: err.Error()}
		respondDeactivateUser(output, m)
		return
	}

	output.OK = true
	respondDeactivateUser(output, m)
}

func respondDeactivateUser(output DeactivateUserOutput, m *nats.Msg) {
	bs, _ := json.Marshal(output)
	m.Respond(bs)
}
