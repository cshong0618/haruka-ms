package auth

import (
	"context"
	"errors"
	"time"
)

type AuthenticationRepository interface {
	CreateAuthenticationStore(ctx context.Context, userID string, handler string, password string) error
	FindByHandler(ctx context.Context, handler string) (UserAuthenticationStore, error)
}

type PasswordHandler interface {
	EncryptPassword(ctx context.Context, rawPassword string) string
	ComparePassword(ctx context.Context, inputPassword string, currentPassword string) bool
}

type TokenRepository interface {
	StoreToken(ctx context.Context, token string, duration time.Duration) error
	CheckToken(ctx context.Context, token string) error
	GenerateToken(ctx context.Context, userID string, duration time.Duration) (string, error)
}

type EventPublisher interface {
	PublishNewAuthSuccessEvent(ctx context.Context, userID string)
	PublishNewAuthFailedEvent(ctx context.Context, userID string)
}

var (
	AuthenticationNotFound = errors.New("handler does not exist")
	HandlerAlreadyExists   = errors.New("handler already exists")
	PasswordNotMatch       = errors.New("password does not match")
	PasswordNotExist       = errors.New("no password")
)
