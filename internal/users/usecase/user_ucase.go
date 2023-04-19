package usecase

import (
	"context"
	"fmt"

	"refactoring/internal/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func (u *userUsecase) Create(ctx context.Context, displayName, email string) (domain.User, error) {
	user, err := u.userRepo.Create(ctx, displayName, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("error create: %s", err.Error())
	}
	return user, nil
}

func (u *userUsecase) GetByID(ctx context.Context, userID string) (domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("error get user by ID: %s", err.Error())
	}
	return user, nil
}

func (u *userUsecase) GetAll(ctx context.Context) (domain.UserList, error) {
	list, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error with get list err: %s", err.Error())
	}
	return list, nil
}

func (u *userUsecase) Update(ctx context.Context, userID, displayName string) error {
	err := u.userRepo.Update(ctx, userID, displayName)
	if err != nil {
		return fmt.Errorf("error with update user: %s", err.Error())
	}
	return nil
}
func (u *userUsecase) Delete(ctx context.Context, userID string) error {
	err := u.userRepo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("error with delete user %s", err.Error())
	}
	return nil
}
