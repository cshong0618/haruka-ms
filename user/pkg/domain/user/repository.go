package user

import "context"

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	UpdateStatus(ctx context.Context, ID string, status Status) (User, error)
	FindById(ctx context.Context, ID string) (User, error)
	FindAll(ctx context.Context) ([]User, error)
}
