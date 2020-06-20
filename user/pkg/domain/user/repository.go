package user

import "context"

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	FindById(ctx context.Context, ID string) (User, error)
	FindAll(ctx context.Context) ([]User, error)
}
