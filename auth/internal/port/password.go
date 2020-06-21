package port

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHandler struct{}

func NewPasswordHandler() *PasswordHandler {
	return &PasswordHandler{}
}

func (p PasswordHandler) EncryptPassword(ctx context.Context, rawPassword string) string {
	bs, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 10)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (p PasswordHandler) ComparePassword(ctx context.Context, inputPassword string, currentPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(inputPassword))
	return err == nil
}
