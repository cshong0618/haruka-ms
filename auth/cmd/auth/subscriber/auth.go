package subscriber

import (
	"context"
	"encoding/json"
	"github.com/cshong0618/haruka/auth/pkg/usecase"
	"github.com/nats-io/nats.go"
)

type AuthenticationSubscriber struct {
	authService *usecase.AuthenticationService
}

func NewAuthenticationSubscriber(authService *usecase.AuthenticationService) *AuthenticationSubscriber {
	return &AuthenticationSubscriber{authService: authService}
}

type AuthenticationResult struct {
	OK bool `json:"ok"`
}

type CreateAuthenticationInput struct {
	UserID   string `json:"userId"`
	Handler  string `json:"handler"`
	Password string `json:"password"`
}

type CreateAuthenticationOutput struct {
	AuthenticationResult
	Error *Error `json:"error"`
}

func (subscriber *AuthenticationSubscriber) CreateAuthentication(m *nats.Msg) {
	var input CreateAuthenticationInput
	err := json.Unmarshal(m.Data, &input)
	output := CreateAuthenticationOutput{}

	if err != nil {
		output.OK = false
		output.Error = &Error{Reason: err.Error()}
		respondCreateAuthentication(output, m)
		return
	}

	err = subscriber.authService.CreateAuthenticationStore(context.Background(),
		input.UserID,
		input.Handler,
		input.Password)

	if err != nil {
		output.OK = false
		output.Error = &Error{Reason: err.Error()}
		respondCreateAuthentication(output, m)
		return
	}

	output.OK = true
	respondCreateAuthentication(output, m)
}

func respondCreateAuthentication(output CreateAuthenticationOutput, m *nats.Msg) {
	bs, _ := json.Marshal(output)
	m.Respond(bs)
}
