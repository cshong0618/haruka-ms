package repository

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"time"
)

const secret = "auth-secret-example"

type TokenRepository struct{
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (t *TokenRepository) StoreToken(ctx context.Context, token string, duration time.Duration) error {
	cmd := t.client.Set(ctx, token, "", duration)
	return cmd.Err()
}

func (t *TokenRepository) CheckToken(ctx context.Context, token string) error {
	cmd := t.client.Get(ctx, token)
	return cmd.Err()
}

func (t *TokenRepository) GenerateToken(ctx context.Context, userID string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp": time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	return tokenString, err
}
