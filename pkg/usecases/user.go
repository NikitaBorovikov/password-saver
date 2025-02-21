package usecases

import (
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
)

type UserUseCase struct {
	UserRepository model.UserRepository
}

func NewUserUseCase(ur model.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepository: ur,
	}
}

func (uc *UserUseCase) Registration(u *model.User) error {
	return nil
}

func (uc *UserUseCase) LogIn(q *dto.LogInRequest) (*dto.UserResponse, error) {
	return nil, nil
}

func (uc *UserUseCase) Update(u *model.User) error {
	return nil
}

func (uc *UserUseCase) Delete(userID int64) error {
	return nil
}
