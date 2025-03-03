package usecases

import (
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
)

type PasswordUseCase struct {
	PasswordRepository model.PasswordRepository
}

func NewPasswordUseCase(pr model.PasswordRepository) *PasswordUseCase {
	return &PasswordUseCase{
		PasswordRepository: pr,
	}
}

func (uc *PasswordUseCase) Save(req *dto.SavePasswordRequest) error {
	return nil
}

func (uc *PasswordUseCase) GetAll(userID int64) ([]model.Password, error) {
	return nil, nil
}

func (uc *PasswordUseCase) GetByID(passwordID string) (*model.Password, error) {
	return nil, nil
}

func (uc *PasswordUseCase) Update(p *model.Password) error {
	return nil
}

func (uc *PasswordUseCase) Delete(passwordID string) error {
	return nil
}
