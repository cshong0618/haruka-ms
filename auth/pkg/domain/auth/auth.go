package auth

import "time"

type UserAuthenticationStore struct {
	UserID    string
	Handler   string
	Passwords []Password
}

type Password struct {
	Value     string
	CreatedOn time.Time
}
