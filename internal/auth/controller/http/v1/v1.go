package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"project/internal/auth/domain"
	"project/internal/auth/service"
)

type app interface {
	Register(ctx context.Context, req service.RegisterParams) (service.RegisterResult, error)
	Login(ctx context.Context, req service.LoginParams) (service.LoginResult, error)
	Logout(ctx context.Context, req service.LogoutParams) error
}

type V1 struct {
	uc       app
	validate *validator.Validate
}

func New(uc app) *V1 {
	return &V1{
		uc:       uc,
		validate: validator.New(),
	}
}

func (v V1) Validate(req any) error {
	err := v.validate.Struct(req)
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err
	}

	return ValidationError(ve[0])
}

func ValidationError(fe validator.FieldError) error {
	switch fe.Field() {
	case "Email":
		switch fe.Tag() {
		case "required":
			return domain.ErrEmailEmpty
		case "email":
			return service.ErrInvalidCredentials
		}
	case "Password":
		switch fe.Tag() {
		case "required":
			return domain.ErrPasswordEmpty
		case "min":
			return domain.ErrPasswordShort
		case "max":
			return domain.ErrPasswordLong
		}
	}
	return fmt.Errorf("validation failed: %s", fe.Field())
}
