package usecase

import (
	"context"
	"github.com/cshong0618/haruka/user/pkg/domain/user"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository user.Repository
}

func NewUserService(userRepository user.Repository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (service *UserService) CreateUser(
	ctx context.Context,
	name string) (user.User, error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Info("create user")

	newUser := user.User{
		Name: name,
	}

	newUser, err := service.userRepository.Create(ctx, newUser)
	if err != nil {
		return user.User{}, err
	}

	return newUser, nil
}

func (service *UserService) FindUserById(
	ctx context.Context,
	ID string) (user.User, error) {
	log.WithFields(log.Fields{
		"id": ID,
	}).Info("find user by id")

	targetUser, err := service.userRepository.FindById(ctx, ID)
	if err != nil {
		return user.User{}, err
	}

	return targetUser, nil
}

func (service *UserService) FindAll(ctx context.Context) ([]user.User, error) {
	log.Info("find all users")

	users, err := service.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}