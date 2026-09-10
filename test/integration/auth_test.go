//go:build integration

package integration_test

import (
	"github.com/google/uuid"

	"project/pkg/httpclient/authclient"
)

func (s *Suite) Test_Auth_RegisterLoginLogout() {
	email := s.uniqueEmail()

	userID, err := s.auth.Register(ctx, email, testPassword)
	s.NoError(err)
	s.NotEqual(uuid.Nil, userID)

	token, err := s.auth.Login(ctx, email, testPassword)
	s.NoError(err)
	s.NotEmpty(token)

	s.auth.SetToken(token)
	s.NoError(s.auth.Logout(ctx))

	// повторный logout тем же токеном — токен уже отозван
	err = s.auth.Logout(ctx)
	s.ErrorIs(err, authclient.ErrUnauthorized)
}

func (s *Suite) Test_Auth_RegisterDuplicateEmail() {
	email := s.uniqueEmail()

	_, err := s.auth.Register(ctx, email, testPassword)
	s.NoError(err)

	_, err = s.auth.Register(ctx, email, testPassword)
	s.ErrorIs(err, authclient.ErrEmailAlreadyExists)
}

func (s *Suite) Test_Auth_LoginWrongPassword() {
	email := s.uniqueEmail()

	_, err := s.auth.Register(ctx, email, testPassword)
	s.NoError(err)

	_, err = s.auth.Login(ctx, email, "wrong-password")
	s.ErrorIs(err, authclient.ErrInvalidCredentials)
}

func (s *Suite) Test_Auth_LoginUnknownEmail() {
	_, err := s.auth.Login(ctx, s.uniqueEmail(), testPassword)
	s.ErrorIs(err, authclient.ErrInvalidCredentials)
}

func (s *Suite) Test_Auth_RegisterEmptyPassword() {
	_, err := s.auth.Register(ctx, s.uniqueEmail(), "")
	s.ErrorIs(err, authclient.ErrValidation)
}

func (s *Suite) Test_Auth_LogoutWithoutToken() {
	client := authclient.New(authclient.Config{Host: "localhost", Port: "8081"}) // токен не установлен

	err := client.Logout(ctx)
	s.ErrorIs(err, authclient.ErrUnauthorized)
}
