package usecase

import (
	"context"
	"github.com/cshong0618/haruka/auth/pkg/domain/auth"
	log "github.com/sirupsen/logrus"
	"time"
)

var tokenValidity = 24 * time.Hour

type AuthenticationService struct {
	authRepository  auth.AuthenticationRepository
	passwordHandler auth.PasswordHandler
	tokenRepository auth.TokenRepository
	eventPublisher  auth.EventPublisher
}

func NewAuthenticationService(
	authRepository auth.AuthenticationRepository,
	passwordHandler auth.PasswordHandler,
	tokenRepository auth.TokenRepository,
	eventPublisher auth.EventPublisher) *AuthenticationService {
	return &AuthenticationService{
		authRepository: authRepository,
		passwordHandler: passwordHandler,
		tokenRepository: tokenRepository,
		eventPublisher: eventPublisher,
	}
}

func (service *AuthenticationService) CreateAuthenticationStore(
	ctx context.Context,
	userID string,
	handle string,
	password string) error {
	encryptedPassword := service.passwordHandler.EncryptPassword(ctx, password)
	log.WithFields(log.Fields{
		"userId": userID,
		"handle": handle,
	}).Info("create authentication")

	err := service.authRepository.CreateAuthenticationStore(ctx, userID, handle, encryptedPassword)
	if err != nil {
		log.WithFields(log.Fields{
			"handle": handle,
			"userId": userID,
		}).Info("create authentication failed")
		service.eventPublisher.PublishNewAuthFailedEvent(ctx, userID)
	} else {
		log.WithFields(log.Fields{
			"handle": handle,
			"userId": userID,
		}).Info("create authentication success")
		service.eventPublisher.PublishNewAuthSuccessEvent(ctx, userID)
	}

	return err
}

func (service *AuthenticationService) Login(
	ctx context.Context,
	handler string,
	password string) error {
	log.WithFields(log.Fields{
		"handle": handler,
	}).Info("login")

	authenticationStore, err := service.authRepository.FindByHandler(ctx, handler)
	if err != nil {
		logLoginError(handler, err)
		return err
	}

	if len(authenticationStore.Passwords) == 0 {
		err := auth.PasswordNotExist
		logLoginError(handler, err)
		return err
	}

	currentPassword := loginGetLatestPassword(authenticationStore.Passwords)
	if currentPassword == "" {
		err := auth.PasswordNotExist
		logLoginError(handler, err)
		return err
	}

	ok := service.passwordHandler.ComparePassword(ctx, password, currentPassword)
	if !ok {
		err := auth.PasswordNotMatch
		logLoginError(handler, err)
		return err
	}

	token, err := service.tokenRepository.GenerateToken(ctx, authenticationStore.UserID, 0)
	if err != nil {
		logLoginError(handler, err)
		return err
	}

	err = service.tokenRepository.StoreToken(ctx, token, tokenValidity)
	if err != nil {
		logLoginError(handler, err)
		return err
	}

	return nil
}

func logLoginError(handler string, err error) {
	log.WithFields(log.Fields{
		"handle": handler,
		"error":   err.Error(),
	}).Info("login failed")
}

func loginGetLatestPassword(passwords []auth.Password) string {
	genesis := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	currentPassword := ""

	for _, password := range passwords {
		if password.CreatedOn.After(genesis) {
			currentPassword = password.Value
		}
	}

	return currentPassword
}
