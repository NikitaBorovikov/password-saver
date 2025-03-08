package usecases

import (
	"fmt"
	"password-saver/pkg/config"
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
	"password-saver/pkg/usecases/encryption"

	"github.com/go-playground/validator"
)

type PasswordUseCase struct {
	PasswordRepository model.PasswordRepository
	cfg                *config.EncryptKeys
}

func NewPasswordUseCase(pr model.PasswordRepository, cfg *config.EncryptKeys) *PasswordUseCase {
	return &PasswordUseCase{
		PasswordRepository: pr,
		cfg:                cfg,
	}
}

func (uc *PasswordUseCase) Save(req *dto.SavePasswordRequest) error {

	if err := validateForSavePassword(req); err != nil {
		return err
	}

	plainPassword := []byte(req.Password)
	plainService := []byte(req.Service)

	encPassword, err := encryption.Encrypt(plainPassword, []byte(uc.cfg.EncPasswordKey))
	if err != nil {
		return err
	}

	encService, err := encryption.Encrypt(plainService, []byte(uc.cfg.EncServiceKey))
	if err != nil {
		return err
	}

	p := &model.Password{
		EncService:  encService,
		EncPassword: encPassword,
	}

	if err := uc.PasswordRepository.Save(p); err != nil {
		return err
	}
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

func validateForSavePassword(req *dto.SavePasswordRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("failed to validate password struct: %v", err)
	}
	return nil
}
