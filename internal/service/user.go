package service

import (
	"context"

	"tmpl/internal/model"
	"tmpl/internal/repository"
)

type User struct {
	userRepo *repository.User
}

func NewUser(userRepo *repository.User) *User {
	return &User{
		userRepo: userRepo,
	}
}

func (u *User) GetOrCreateUser(
	ctx context.Context,
	email string,
	googleID string,
) (*model.User, error) {
	user, err := u.userRepo.GetOrCreateUserWithGoogleID(ctx, email, googleID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
